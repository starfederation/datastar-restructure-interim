import { Datastar } from "library/engine";
import { RemoteActionPlugin } from "library/plugins/actions/backend/remote";
import {
    DeleteSSEActionPlugin,
    GetSSEActionPlugin,
    PatchSSEActionPlugin,
    PostSSEActionPlugin,
    PutSSEActionPlugin,
} from "library/plugins/actions/backend/sse";
import { ClipboardActionPlugin } from "library/plugins/actions/dom/clipboard";
import { SetAllActionPlugin } from "library/plugins/actions/logic/setAll";
import { ToggleAllActionPlugin } from "library/plugins/actions/logic/toggleAll";
import { ClampFitActionPlugin } from "library/plugins/actions/math/clampFit";
import { ClampFitIntActionPlugin } from "library/plugins/actions/math/clampFitInt";
import { FitActionPlugin } from "library/plugins/actions/math/fit";
import { FitIntActionPlugin } from "library/plugins/actions/math/fitInt";
import { ScrollActionPlugin } from "library/plugins/actions/visibility/scroll";
import { FetchIndicatorPlugin } from "library/plugins/attributes/backend/fetch-indicator";
import { HeadersPlugin } from "library/plugins/attributes/backend/headers";
import { BindAttributePlugin } from "library/plugins/attributes/basics/bind";
import { ClassAttributePlugin } from "library/plugins/attributes/basics/class";
import { ModelAttributePlugin } from "library/plugins/attributes/basics/model";
import { OnAttributePlugin } from "library/plugins/attributes/basics/on";
import { TextAttributePlugin } from "library/plugins/attributes/basics/text";
import { IntersectionAttributePlugin } from "library/plugins/attributes/visibility/intersects";
import { ScrollIntoViewAttributePlugin } from "library/plugins/attributes/visibility/scrollIntoView";
import { ShowAttributePlugin } from "library/plugins/attributes/visibility/show";
import { TeleportAttributePlugin } from "library/plugins/attributes/visibility/teleport";
import { ViewTransitionAttributePlugin } from "library/plugins/attributes/visibility/viewTransition";

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
