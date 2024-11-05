import { PreprocessorPlugin, RegexpGroups } from "library/src/engine";
import { wholePrefixSuffix } from "library/src/utils/regex";

// Replacing ~foo with ctx.refs.foo
export const RefProcessorPlugin: PreprocessorPlugin = {
    name: "ref",
    pluginType: "preprocessor",
    regexp: wholePrefixSuffix("~", "ref", "", false),
    replacer({ ref }: RegexpGroups) {
        return `document.querySelector(ctx.store()._dsPlugins.refs.${ref})`;
    },
};
