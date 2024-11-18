import { PreprocessorPlugin, RegexpGroups } from "../../../../engine";
import { PLUGIN_PREPROCESSOR } from "../../../../engine/client_only_consts";
import { wholePrefixSuffix } from "../../../../utils/regex";

// Replacing $signal with ctx.store.signal.value`
export const SignalsProcessor: PreprocessorPlugin = {
    name: "signal",
    pluginType: PLUGIN_PREPROCESSOR,
    regexp: wholePrefixSuffix("\\$", "signal", "(?<method>\\([^\\)]*\\))?"),
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
