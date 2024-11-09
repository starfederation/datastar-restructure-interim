import { PreprocessorPlugin, RegexpGroups } from "library/src/engine";
import { wholePrefixSuffix } from "library/src/utils/regex";

// Replacing $signal with ctx.store.signal.value`
export const SignalsProcessorPlugin: PreprocessorPlugin = {
    name: "signal",
    pluginType: "preprocessor",
    regexp: wholePrefixSuffix("$", "signal", "(?<method>\\([^\\)]*\\))?"),
    replacer: (groups: RegexpGroups) => {
        const { signal, method } = groups;
        const prefix = `ctx.store()`;
        if (!method?.length) {
            return `${prefix}.${signal}.value`;
        }
        const parts = signal.split(".");
        const methodName = parts.pop();
        const nestedSignal = parts.join(".");
        return `${prefix}.${nestedSignal}.value.${methodName}${method}`;
    },
};
