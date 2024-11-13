import { Ref } from "../plugins/attributes/core/ref";
import { Store } from "../plugins/attributes/core/store";
import { ActionsProcessor } from "../plugins/preprocessors/core/actions";
import { RefProcessor } from "../plugins/preprocessors/core/ref";
import { SignalsProcessor } from "../plugins/preprocessors/core/signals";
import { Engine } from "./engine";

export * from "./const";
export type * from "./types";

const ds = new Engine();
ds.load(
    ActionsProcessor,
    SignalsProcessor,
    RefProcessor,
    Store,
    Ref,
);

export const Datastar = ds;
