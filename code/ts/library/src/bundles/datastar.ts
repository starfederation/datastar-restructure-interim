import { Datastar } from "library/src/engine";
import { RemoteActionPlugin } from "library/src/plugins/actions/backend/remote";
import {
    DeleteSSEActionPlugin,
    GetSSEActionPlugin,
    PatchSSEActionPlugin,
    PostSSEActionPlugin,
    PutSSEActionPlugin,
} from "library/src/plugins/actions/backend/sse";
import { ClipboardActionPlugin } from "library/src/plugins/actions/dom/clipboard";
import { SetAllActionPlugin } from "library/src/plugins/actions/logic/setAll";
import { ToggleAllActionPlugin } from "library/src/plugins/actions/logic/toggleAll";
import { ClampFitActionPlugin } from "library/src/plugins/actions/math/clampFit";
import { ClampFitIntActionPlugin } from "library/src/plugins/actions/math/clampFitInt";
import { FitActionPlugin } from "library/src/plugins/actions/math/fit";
import { FitIntActionPlugin } from "library/src/plugins/actions/math/fitInt";
import { ScrollActionPlugin } from "library/src/plugins/actions/visibility/scroll";
import { FetchIndicatorPlugin } from "library/src/plugins/attributes/backend/fetch-indicator";
import { HeadersPlugin } from "library/src/plugins/attributes/backend/headers";
import { BindAttributePlugin } from "library/src/plugins/attributes/basics/bind";
import { ClassAttributePlugin } from "library/src/plugins/attributes/basics/class";
import { ModelAttributePlugin } from "library/src/plugins/attributes/basics/model";
import { OnAttributePlugin } from "library/src/plugins/attributes/basics/on";
import { TextAttributePlugin } from "library/src/plugins/attributes/basics/text";
import { IntersectionAttributePlugin } from "library/src/plugins/attributes/visibility/intersects";
import { ScrollIntoViewAttributePlugin } from "library/src/plugins/attributes/visibility/scrollIntoView";
import { ShowAttributePlugin } from "library/src/plugins/attributes/visibility/show";
import { TeleportAttributePlugin } from "library/src/plugins/attributes/visibility/teleport";
import { ViewTransitionAttributePlugin } from "library/src/plugins/attributes/visibility/viewTransition";

Datastar.load(
    RemoteActionPlugin,
    GetSSEActionPlugin,
    PostSSEActionPlugin,
    PutSSEActionPlugin,
    PatchSSEActionPlugin,
    DeleteSSEActionPlugin,
    ClipboardActionPlugin,
    SetAllActionPlugin,
    ToggleAllActionPlugin,
    FitActionPlugin,
    FitIntActionPlugin,
    ClampFitActionPlugin,
    ClampFitIntActionPlugin,
    ScrollActionPlugin,
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
