import { Ref } from "../plugins/attributes/core/ref";
import { StoreAttributePlugin } from "../plugins/attributes/core/store";
import { ActionsProcessorPlugin } from "../plugins/preprocessors/core/actions";
import { RefProcessor } from "../plugins/preprocessors/core/ref";
import { SignalsProcessorPlugin } from "../plugins/preprocessors/core/signals";
import { Engine } from "./engine";

export * from "./const";
export type * from "./types";

const ds = new Engine();
ds.load(
    ActionsProcessorPlugin,
    SignalsProcessorPlugin,
    RefProcessor,
    StoreAttributePlugin,
    Ref,
);

export const Datastar = ds;
