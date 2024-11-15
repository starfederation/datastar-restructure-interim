import { Store } from "../plugins/official/attributes/core/store";
import { ActionsProcessor } from "../plugins/official/preprocessors/core/actions";
import { SignalsProcessor } from "../plugins/official/preprocessors/core/signals";
import { Engine } from "./engine";

export * from "./const";
export type * from "./types";

const ds = new Engine();
ds.load(
    ActionsProcessor,
    SignalsProcessor,
    Store,
);

export const Datastar = ds;
