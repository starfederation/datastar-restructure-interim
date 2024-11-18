import { PreprocessorPlugin, RegexpGroups } from "../../../../engine";
import { PLUGIN_PREPROCESSOR } from "../../../../engine/client_only_consts";
import { wholePrefixSuffix } from "../../../../utils/regex";

// Replacing $action(args) with ctx.actions.action(ctx, args)
export const ActionsProcessor: PreprocessorPlugin = {
    name: "action",
    pluginType: PLUGIN_PREPROCESSOR,
    regexp: wholePrefixSuffix(
        "\\$",
        "action",
        "(?<call>\\((?<args>.*)\\))",
        false,
    ),
    replacer: ({ action, args }: RegexpGroups) => {
        const withCtx = [`ctx`];
        if (args) {
            withCtx.push(...args.split(",").map((x) => x.trim()));
        }
        const argsJoined = withCtx.join(",");
        return `ctx.actions.${action}.method(${argsJoined})`;
    },
};
