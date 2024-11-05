import { AttributePlugin, RegexpGroups } from "library/src/engine";

const dataURIRegex = /^data:(?<mime>[^;]+);base64,(?<contents>.*)$/;
const updateModelEvents = ["change", "input", "keydown"];

export const ModelAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "model",
    mustHaveEmptyKey: true,
    preprocessors: {
        post: [
            {
                pluginType: "preprocessor",
                name: "model",
                regexp: /(?<whole>.+)/g,
                replacer: (groups: RegexpGroups) => {
                    const { whole } = groups;
                    return `ctx.store().${whole}`;
                },
            },
        ],
    },
    // bypassExpressionFunctionCreation: () => true,
    onLoad: (ctx) => {
        const { el, expression } = ctx;
        const signal = ctx.expressionFn(ctx);
        const tnl = el.tagName.toLowerCase();

        if (expression.startsWith("ctx.store().ctx.store()")) {
            throw new Error(
                `Model attribute on #${el.id} must have a signal name, you probably prefixed with $ by accident`,
            );
        }

        const isInput = tnl.includes("input");
        const type = el.getAttribute("type");
        const isCheckbox = tnl.includes("checkbox") ||
            (isInput && type === "checkbox");
        const isSelect = tnl.includes("select");
        const isRadio = tnl.includes("radio") || (isInput && type === "radio");
        const isFile = isInput && type === "file";

        const signalName = expression.replaceAll("ctx.store().", "");
        if (isRadio) {
            const name = el.getAttribute("name");
            if (!name?.length) {
                el.setAttribute("name", signalName);
            }
        }

        const setInputFromSignal = () => {
            if (!signal) {
                throw new Error(`Signal ${signalName} not found`);
            }
            const hasValue = "value" in el;
            const v = signal.value;
            if (isCheckbox || isRadio) {
                const input = el as HTMLInputElement;
                if (isCheckbox) {
                    input.checked = v;
                } else if (isRadio) {
                    // evaluate the value as string to handle any type casting
                    // automatically since the attribute has to be a string anyways
                    input.checked = `${v}` === input.value;
                }
            } else if (isFile) {
                // File input reading from a signal is not supported yet
            } else if (isSelect) {
                const select = el as HTMLSelectElement;
                if (select.multiple) {
                    const v = signal.value;
                    Array.from(select.options).forEach((opt) => {
                        if (opt?.disabled) return;
                        opt.selected = v.includes(opt.value);
                    });
                } else {
                    select.value = `${v}`;
                }
            } else if (hasValue) {
                el.value = `${v}`;
            } else {
                el.setAttribute("value", `${v}`);
            }
        };
        const cleanupSetInputFromSignal = ctx.reactivity.effect(
            setInputFromSignal,
        );

        const setSignalFromInput = async () => {
            if (isFile) {
                const files = [...((el as HTMLInputElement)?.files || [])],
                    allContents: string[] = [],
                    allMimes: string[] = [],
                    allNames: string[] = [];

                await Promise.all(
                    files.map((f) => {
                        return new Promise<void>((resolve) => {
                            const reader = new FileReader();
                            reader.onload = () => {
                                if (typeof reader.result !== "string") {
                                    throw new Error(
                                        `Invalid result type: ${typeof reader
                                            .result}`,
                                    );
                                }
                                const match = reader.result.match(dataURIRegex);
                                if (!match?.groups) {
                                    throw new Error(
                                        `Invalid data URI: ${reader.result}`,
                                    );
                                }
                                allContents.push(match.groups.contents);
                                allMimes.push(match.groups.mime);
                                allNames.push(f.name);
                            };
                            reader.onloadend = () => resolve(void 0);
                            reader.readAsDataURL(f);
                        });
                    }),
                );

                signal.value = allContents;
                const s = ctx.store();
                const mimeName = `${signalName}Mimes`,
                    nameName = `${signalName}Names`;
                if (mimeName in s) {
                    s[`${mimeName}`].value = allMimes;
                }
                if (nameName in s) {
                    s[`${nameName}`].value = allNames;
                }
                return;
            }

            const current = signal.value;
            const input = (el as HTMLInputElement) || (el as HTMLElement);

            if (typeof current === "number") {
                signal.value = Number(
                    input.value || input.getAttribute("value"),
                );
            } else if (typeof current === "string") {
                signal.value = input.value || input.getAttribute("value") || "";
            } else if (typeof current === "boolean") {
                if (isCheckbox) {
                    signal.value = input.checked ||
                        input.getAttribute("checked") === "true";
                } else {
                    signal.value = Boolean(
                        input.value || input.getAttribute("value"),
                    );
                }
            } else if (typeof current === "undefined") {
            } else if (typeof current === "bigint") {
                signal.value = BigInt(
                    input.value || input.getAttribute("value") || "0",
                );
            } else if (Array.isArray(current)) {
                // check if the input is a select element
                if (isSelect) {
                    const select = el as HTMLSelectElement;
                    const selectedOptions = [...select.selectedOptions];
                    const selectedValues = selectedOptions.map((opt) =>
                        opt.value
                    );
                    signal.value = selectedValues;
                } else {
                    signal.value = JSON.parse(input.value).split(",");
                }
                console.log(input.value);
            } else {
                console.log(typeof current);
                throw new Error(
                    `Unsupported type ${typeof current} for signal ${signalName}`,
                );
            }
        };

        const parts = el.tagName.split("-");
        const isCustomElement = parts.length > 1;
        if (isCustomElement) {
            const customElementPrefix = parts[0].toLowerCase();
            updateModelEvents.forEach((eventType) => {
                updateModelEvents.push(`${customElementPrefix}-${eventType}`);
            });
        }

        updateModelEvents.forEach((eventType) =>
            el.addEventListener(eventType, setSignalFromInput)
        );

        return () => {
            cleanupSetInputFromSignal();
            updateModelEvents.forEach((event) =>
                el.removeEventListener(event, setSignalFromInput)
            );
        };
    },
};
