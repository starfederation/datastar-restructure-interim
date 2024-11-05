import { RefProcessorPlugin } from "library/src/plugins/actions/core/ref";
import { ComputedPlugin } from "library/src/plugins/attributes/core/computed";
import { RefPlugin } from "library/src/plugins/attributes/core/ref";
import { StoreAttributePlugin } from "library/src/plugins/attributes/core/store";
import { ActionsProcessorPlugin } from "library/src/plugins/preprocessors/actions";
import { SignalsProcessorPlugin } from "library/src/plugins/preprocessors/signals";
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
