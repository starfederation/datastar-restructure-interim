import { RefProcessorPlugin } from "../plugins/actions/core/ref";
import { ComputedPlugin } from "../plugins/attributes/core/computed";
import { RefPlugin } from "../plugins/attributes/core/ref";
import { StoreAttributePlugin } from "../plugins/attributes/core/store";
import { ActionsProcessorPlugin } from "../plugins/preprocessors/core/actions";
import { SignalsProcessorPlugin } from "../plugins/preprocessors/core/signals";
import { Engine } from "./engine";

export * from "./const";
export type * from "./types";

const ds = new Engine();
ds.load(
    ActionsProcessorPlugin,
    SignalsProcessorPlugin,
    RefProcessorPlugin,
    StoreAttributePlugin,
    ComputedPlugin,
    RefPlugin,
);

export const Datastar = ds;
