import { RefProcessorPlugin } from "library/plugins/actions/core/ref";
import { ComputedPlugin } from "library/plugins/attributes/core/computed";
import { RefPlugin } from "library/plugins/attributes/core/ref";
import { StoreAttributePlugin } from "library/plugins/attributes/core/store";
import { ActionsProcessorPlugin } from "library/plugins/preprocessors/actions";
import { SignalsProcessorPlugin } from "library/plugins/preprocessors/signals";
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
