import { Datastar } from "../engine";
import { IsFetching } from "../plugins/official/actions/backend/isFetching";
import { RemoteSignals } from "../plugins/official/actions/backend/remote";
import { DeleteSSE } from "../plugins/official/actions/backend/sseDelete";
import { GetSSE } from "../plugins/official/actions/backend/sseGet";
import { PatchSSE } from "../plugins/official/actions/backend/ssePatch";
import { PostSSE } from "../plugins/official/actions/backend/ssePost";
import { PutSSE } from "../plugins/official/actions/backend/ssePut";
import { Clipboard } from "../plugins/official/actions/dom/clipboard";
import { RefAction } from "../plugins/official/actions/dom/ref";
import { SetAll } from "../plugins/official/actions/logic/setAll";
import { ToggleAll } from "../plugins/official/actions/logic/toggleAll";
import { ClampFit } from "../plugins/official/actions/math/clampFit";
import { ClampFitInt } from "../plugins/official/actions/math/clampFitInt";
import { Fit } from "../plugins/official/actions/math/fit";
import { FitInt } from "../plugins/official/actions/math/fitInt";
import { FetchIndicator } from "../plugins/official/attributes/backend/fetchIndicator";
import { Header } from "../plugins/official/attributes/backend/header";
import { ReplaceUrl } from "../plugins/official/attributes/backend/replaceUrl";
import { Bind } from "../plugins/official/attributes/dom/bind";
import { Class } from "../plugins/official/attributes/dom/class";
import { Model } from "../plugins/official/attributes/dom/model";
import { On } from "../plugins/official/attributes/dom/on";
import { RefAttribute } from "../plugins/official/attributes/dom/ref";
import { Text } from "../plugins/official/attributes/dom/text";
import { Persist } from "../plugins/official/attributes/storage/persist";
import { Intersection } from "../plugins/official/attributes/visibility/intersects";
import { ScrollIntoView } from "../plugins/official/attributes/visibility/scrollIntoView";
import { Show } from "../plugins/official/attributes/visibility/show";
import { Teleport } from "../plugins/official/attributes/visibility/teleport";
import { ViewTransition } from "../plugins/official/attributes/visibility/viewTransition";
import { ExecuteScript } from "../plugins/official/watchers/backend/sseExecuteScript";
import { MergeFragments } from "../plugins/official/watchers/backend/sseMergeFragment";
import { MergeSignals } from "../plugins/official/watchers/backend/sseMergeSignals";
import { RemoveFragments } from "../plugins/official/watchers/backend/sseRemoveFragments";
import { RemoveSignals } from "../plugins/official/watchers/backend/sseRemoveSignals";

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
    RefAction,
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
    // attributes/dom
    Bind,
    Class,
    Model,
    On,
    RefAttribute,
    Text,
    // attributes/storage
    Persist,
    // attributes/visibility
    Intersection,
    ScrollIntoView,
    Show,
    Teleport,
    ViewTransition,
    // effects
    MergeFragments,
    MergeSignals,
    RemoveFragments,
    RemoveSignals,
    ExecuteScript,
);
