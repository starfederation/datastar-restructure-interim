import { Datastar } from "../engine";
import { IsFetching } from "../plugins/actions/backend/isFetching";
import { RemoteSignals } from "../plugins/actions/backend/remote";
import { DeleteSSE } from "../plugins/actions/backend/sseDelete";
import { GetSSE } from "../plugins/actions/backend/sseGet";
import { PatchSSE } from "../plugins/actions/backend/ssePatch";
import { PostSSE } from "../plugins/actions/backend/ssePost";
import { PutSSE } from "../plugins/actions/backend/ssePut";
import { Clipboard } from "../plugins/actions/dom/clipboard";
import { SetAll } from "../plugins/actions/logic/setAll";
import { ToggleAll } from "../plugins/actions/logic/toggleAll";
import { ClampFit } from "../plugins/actions/math/clampFit";
import { ClampFitInt } from "../plugins/actions/math/clampFitInt";
import { Fit } from "../plugins/actions/math/fit";
import { FitInt } from "../plugins/actions/math/fitInt";
import { FetchIndicator } from "../plugins/attributes/backend/fetchIndicator";
import { Header } from "../plugins/attributes/backend/header";
import { ReplaceUrl } from "../plugins/attributes/backend/replaceUrl";
import { Bind } from "../plugins/attributes/basics/bind";
import { Class } from "../plugins/attributes/basics/class";
import { Model } from "../plugins/attributes/basics/model";
import { On } from "../plugins/attributes/basics/on";
import { Text } from "../plugins/attributes/basics/text";
import { Persist } from "../plugins/attributes/storage/persist";
import { Intersection } from "../plugins/attributes/visibility/intersects";
import { ScrollIntoView } from "../plugins/attributes/visibility/scrollIntoView";
import { Show } from "../plugins/attributes/visibility/show";
import { Teleport } from "../plugins/attributes/visibility/teleport";
import { ViewTransition } from "../plugins/attributes/visibility/viewTransition";

Datastar.load(
    // actions/backend
    IsFetching,
    RemoteSignals,
    DeleteSSE,
    GetSSE,
    PatchSSE,
    PostSSE,
    PutSSE,
    // actions/dom
    Clipboard,
    // actions/logic
    SetAll,
    ToggleAll,
    // actions/math
    ClampFit,
    ClampFitInt,
    Fit,
    FitInt,
    // attributes/backend
    FetchIndicator,
    Header,
    ReplaceUrl,
    // attributes/basics
    Bind,
    Class,
    Model,
    On,
    Text,
    // attributes/storage
    Persist,
    // attributes/visibility
    Intersection,
    ScrollIntoView,
    Show,
    Teleport,
    ViewTransition,
);
