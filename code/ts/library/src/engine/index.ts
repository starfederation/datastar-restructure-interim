import { Star } from "../plugins/official/attributes/core/advice";
import { Computed } from "../plugins/official/attributes/core/computed";
import { Signals } from "../plugins/official/attributes/core/signals";
import { ActionsProcessor } from "../plugins/official/preprocessors/core/actions";
import { SignalsProcessor } from "../plugins/official/preprocessors/core/signals";
import { Engine } from "./engine";

export * from "./consts";
export type * from "./types";

const ds = new Engine();
ds.load(
    ActionsProcessor,
    SignalsProcessor,
    Signals,
    Computed,
    Star,
);

export const Datastar = ds;
