import { AttributePlugin, DATASTAR } from "library/src/engine";
import { argsToMs } from "../../../utils/arguments";

const DISPLAY = "display";
const NONE = "none";
const IMPORTANT = "important";
const DURATION = "duration";

const SHOW = "show";
const SHOW_CLASS = `${DATASTAR}-showing`;
const HIDE_CLASS = `${DATASTAR}-hiding`;
const SHOW_DURATION_TRANSITION_STYLE = `${DATASTAR}-show-duration-transition`;

export const ShowAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: SHOW,
    "allowedModifiers": new Set([IMPORTANT, DURATION]),

    onLoad: (ctx) => {
        const { el, modifiers, expressionFn, reactivity } = ctx;

        const isImportant = modifiers.has(IMPORTANT);
        const priority = isImportant ? IMPORTANT : undefined;

        let showFn, hideFn;

        const durationArgs = ctx.modifiers.get(DURATION);
        if (durationArgs) {
            let style = document.getElementById(SHOW_DURATION_TRANSITION_STYLE);
            if (!style) {
                style = document.createElement("style");
                style.id = SHOW_DURATION_TRANSITION_STYLE;
                document.head.appendChild(style);
                const durationMs = argsToMs(durationArgs) || "300";
                style.innerHTML = `
            .${SHOW_CLASS} {
              visibility: visible;
              transition: opacity ${durationMs}ms linear;
            }
            .${HIDE_CLASS} {
              visibility: hidden;
              transition: visibility 0s ${durationMs}ms, opacity ${durationMs}ms linear;
            }
          `;
            }

            const transitionEndHandler =
                (classNameToRemove: string) => (event: Event) => {
                    if (event.target === el) {
                        el.classList.remove(classNameToRemove);
                        el.removeEventListener(
                            "transitionend",
                            transitionEndHandler(classNameToRemove),
                        );
                    }
                };

            showFn = () => {
                el.addEventListener(
                    "transitionend",
                    transitionEndHandler(SHOW_CLASS),
                );
                el.classList.add(SHOW_CLASS);
                requestAnimationFrame(() => {
                    el.style.setProperty("opacity", "1", priority);
                });
            };

            hideFn = () => {
                el.addEventListener(
                    "transitionend",
                    transitionEndHandler(HIDE_CLASS),
                );
                el.classList.add(HIDE_CLASS);
                requestAnimationFrame(() => {
                    el.style.setProperty("opacity", "0", priority);
                });
            };
        } else {
            showFn = () => {
                if (el.style.length === 1 && el.style.display === NONE) {
                    el.style.removeProperty(DISPLAY);
                } else {
                    el.style.setProperty(DISPLAY, "", priority);
                }
            };

            hideFn = () => {
                el.style.setProperty(DISPLAY, NONE, priority);
            };
        }

        return reactivity.effect(async () => {
            const expressionEvaluated = await expressionFn(ctx);
            const shouldShow = !!expressionEvaluated;

            if (shouldShow) {
                showFn();
            } else {
                hideFn();
            }
        });
    },
};
