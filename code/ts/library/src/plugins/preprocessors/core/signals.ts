import { PreprocessorPlugin, RegexpGroups } from "../../../engine";
import { wholePrefixSuffix } from "../../../utils/regex";

// Replacing $signal with ctx.store.signal.value`
export const SignalsProcessor: PreprocessorPlugin = {
    name: "signal",
    pluginType: "preprocessor",
    regexp: wholePrefixSuffix("\\$", "signal", ""),
    replacer: (groups: RegexpGroups) => {
        const { signal } = groups;
        const prefix = `ctx.store()`;
        return `${prefix}.${signal}.value`;
    },
};
