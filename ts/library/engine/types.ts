import { HTMLorSVGElement } from "../utils/types";
import { DeepState } from "../vendored/deepsignal";
import { ReadonlySignal, Signal } from "../vendored/preact-core";

export type ExpressionFunction = (ctx: AttributeContext, ...args: any) => any;
export type Reactivity = {
  signal: <T>(value: T) => Signal<T>;
  computed: <T>(fn: () => T) => ReadonlySignal<T>;
  effect: (cb: () => void) => OnRemovalFn;
};

export type AttributeContext = {
  store: () => any;
  mergeStore: (store: DeepState) => void;
  upsertIfMissingFromStore: (path: string, value: any) => void;
  removeFromStore: (...paths: string[]) => void;
  applyPlugins: (target: Element) => void;
  walkSignals: (cb: (name: string, signal: Signal<any>) => void) => void;
  cleanupElementRemovals: (el: Element) => void;
  actions: Readonly<ActionPlugins>;
  reactivity: Reactivity;
  el: Readonly<HTMLorSVGElement>;
  key: Readonly<string>;
  rawKey: Readonly<string>;
  rawExpression: Readonly<string>;
  expression: Readonly<string>;
  expressionFn: ExpressionFunction;
  modifiers: Map<string, string[]>;
  sendDatastarEvent: SendDatastarEvent;
};

export type SendDatastarEvent = (
  category: "core" | "plugin",
  subcategory: string,
  type: string,
  target: Element | Document | Window | string,
  message: string,
  opts?: CustomEventInit,
) => void;

export type InitContext = {
  store: any;
  mergeStore: (store: DeepState) => void;
  actions: Readonly<ActionPlugins>;
  reactivity: Reactivity;
};

export type OnRemovalFn = () => void;
export type AttributePlugin = {
  pluginType: "attribute";
  prefix: string; // The prefix of the `data-${prefix}` attribute
  requiredPluginPrefixes?: Iterable<string>; // If not provided, no plugins are required
  onGlobalInit?: (ctx: InitContext) => void; // Called once on registration of the plugin
  onLoad: (ctx: AttributeContext) => OnRemovalFn | void; // Return a function to be called on removal
  allowedModifiers?: Set<string>; // If not provided, all modifiers are allowed
  mustHaveEmptyExpression?: boolean; // The contents of the data-* attribute must be empty
  mustNotEmptyExpression?: boolean; // The contents of the data-* attribute must not be empty
  mustHaveEmptyKey?: boolean; // The key of the data-* attribute must be empty after the prefix
  mustNotEmptyKey?: boolean; // The key of the data-* attribute must not be empty after the prefix
  allowedTagRegexps?: Set<string>; // If not provided, all tags are allowed
  disallowedTags?: Set<string>; // If not provided, no tags are disallowed
  preprocessors?: {
    pre?: PreprocessorPlugin[];
    post?: PreprocessorPlugin[];
  };
  removeNewLines?: boolean; // If true, the expression is not split by commas
  bypassExpressionFunctionCreation?: (ctx: AttributeContext) => boolean; // If true, the expression function is not created
  argumentNames?: Readonly<string[]>; // The names of the arguments passed to the expression function
};

export type RegexpGroups = Record<string, string>;

export type PreprocessorPlugin = {
  pluginType: "preprocessor";
  name: string;
  regexp: RegExp;
  replacer: (groups: RegexpGroups) => string;
};

export type PreprocessorPlugins = Record<string, PreprocessorPlugin>;

export type ActionMethod = (ctx: AttributeContext, ...args: any[]) => any;

export type ActionPlugin = {
  pluginType: "action";
  name: string;
  method: ActionMethod;
};

export type ActionPlugins = Record<string, ActionPlugin>;

export type DatastarPlugin =
  | AttributePlugin
  | ActionPlugin
  | PreprocessorPlugin;

export interface DatastarEvent {
  time: Date;
  category: "core" | "plugin";
  subcategory: string;
  type: string;
  target: string;
  message: string;
}
