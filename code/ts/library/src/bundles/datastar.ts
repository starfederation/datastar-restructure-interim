import { Datastar } from "library/src/engine";
import { isFetchingActionPlugin } from "../plugins/actions/backend/isFetching";
import { RemoteActionPlugin } from "../plugins/actions/backend/remote";
import {
    DeleteSSEActionPlugin,
    GetSSEActionPlugin,
    PatchSSEActionPlugin,
    PostSSEActionPlugin,
    PutSSEActionPlugin,
} from "../plugins/actions/backend/sse";
import { ClipboardActionPlugin } from "../plugins/actions/dom/clipboard";
import { SetAllActionPlugin } from "../plugins/actions/logic/setAll";
import { ToggleAllActionPlugin } from "../plugins/actions/logic/toggleAll";
import { ClampFitActionPlugin } from "../plugins/actions/math/clampFit";
import { ClampFitIntActionPlugin } from "../plugins/actions/math/clampFitInt";
import { FitActionPlugin } from "../plugins/actions/math/fit";
import { FitIntActionPlugin } from "../plugins/actions/math/fitInt";
import { FetchIndicatorPlugin } from "../plugins/attributes/backend/fetch-indicator";
import { HeadersPlugin } from "../plugins/attributes/backend/headers";
import { BindAttributePlugin } from "../plugins/attributes/basics/bind";
import { ClassAttributePlugin } from "../plugins/attributes/basics/class";
import { ModelAttributePlugin } from "../plugins/attributes/basics/model";
import { OnAttributePlugin } from "../plugins/attributes/basics/on";
import { TextAttributePlugin } from "../plugins/attributes/basics/text";
import { IntersectionAttributePlugin } from "../plugins/attributes/visibility/intersects";
import { ScrollIntoViewAttributePlugin } from "../plugins/attributes/visibility/scrollIntoView";
import { ShowAttributePlugin } from "../plugins/attributes/visibility/show";
import { TeleportAttributePlugin } from "../plugins/attributes/visibility/teleport";
import { ViewTransitionAttributePlugin } from "../plugins/attributes/visibility/viewTransition";

Datastar.load(
    RemoteActionPlugin,
    GetSSEActionPlugin,
    PostSSEActionPlugin,
    PutSSEActionPlugin,
    PatchSSEActionPlugin,
    DeleteSSEActionPlugin,
    isFetchingActionPlugin,
    ClipboardActionPlugin,
    SetAllActionPlugin,
    ToggleAllActionPlugin,
    FitActionPlugin,
    FitIntActionPlugin,
    ClampFitActionPlugin,
    ClampFitIntActionPlugin,
    BindAttributePlugin,
    ClassAttributePlugin,
    ModelAttributePlugin,
    OnAttributePlugin,
    TextAttributePlugin,
    FetchIndicatorPlugin,
    HeadersPlugin,
    IntersectionAttributePlugin,
    ScrollIntoViewAttributePlugin,
    ShowAttributePlugin,
    TeleportAttributePlugin,
    ViewTransitionAttributePlugin,
);
