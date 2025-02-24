export type Page = List | Detail | Form;
export type Command =
  | string
  | [string, ...string[]]
  | {
      name: string;
      args?: string[];
      input?: string;
      dir?: string;
    };
export type Request =
  | string
  | {
      /**
       * The URL to request.
       */
      url: string;
      /**
       * The HTTP method to use.
       */
      method?: string;
      /**
       * The headers to send.
       */
      headers?: {
        [k: string]: string;
      };
      /**
       * The body to send.
       */
      body?: string;
    };
export type Expression =
  | string
  | {
      code: string;
      args?: unknown[];
    };
export type Action =
  | {
      /**
       * The type of the action.
       */
      type: "copy";
      /**
       * The title of the action.
       */
      title?: string;
      /**
       * The text to copy.
       */
      text: string;
      /**
       * The inputs to show when the action is run.
       */
      inputs?: Input[];
      /**
       * The key used as a shortcut.
       */
      key?: string;
    }
  | {
      /**
       * The type of the action.
       */
      type: "paste";
      /**
       * The title of the action.
       */
      title?: string;
      /**
       * The text to paste.
       */
      text: string;
      /**
       * The inputs to show when the action is run.
       */
      inputs?: Input[];
      /**
       * The key used as a shortcut.
       */
      key?: string;
    }
  | {
      /**
       * The type of the action.
       */
      type: "eval";
      /**
       * The title of the action.
       */
      title?: string;
      expression: Expression;
      onSuccess?: OnSuccess;
      /**
       * The inputs to show when the action is run.
       */
      inputs?: Input[];
      /**
       * The key used as a shortcut.
       */
      key?: string;
    }
  | {
      /**
       * The type of the action.
       */
      type: "open";
      /**
       * The title of the action.
       */
      title?: string;
      /**
       * The inputs to show when the action is run.
       */
      inputs?: Input[];
      /**
       * The key used as a shortcut.
       */
      key?: string;
      /**
       * The target to open.
       */
      target: string;
    }
  | {
      /**
       * The type of the action.
       */
      type: "reload";
      /**
       * The title of the action.
       */
      title?: string;
      /**
       * The inputs to show when the action is run.
       */
      inputs?: Input[];
      /**
       * The key used as a shortcut.
       */
      key?: string;
      command: Command;
    }
  | {
      /**
       * The type of the action.
       */
      type: "fetch";
      /**
       * The title of the action.
       */
      title?: string;
      /**
       * The inputs to show when the action is run.
       */
      inputs?: Input[];
      /**
       * The key used as a shortcut.
       */
      key?: string;
      request: Request;
      onSuccess?: OnSuccess;
    }
  | {
      /**
       * The type of the action.
       */
      type: "exec";
      /**
       * The title of the action.
       */
      title?: string;
      /**
       * The inputs to show when the action is run.
       */
      inputs?: Input[];
      /**
       * The key used as a shortcut.
       */
      key?: string;
      command: Command;
      onSuccess?: OnSuccess;
    }
  | {
      /**
       * The type of the action.
       */
      type: "push";
      /**
       * The title of the action.
       */
      title?: string;
      /**
       * The key used as a shortcut.
       */
      key?: string;
      /**
       * The inputs to show when the action is run.
       */
      inputs?: Input[];
      command?: Command;
      request?: Request;
      target?: string;
      expression?: Expression;
      page?: string;
    };
export type Input =
  | {
      /**
       * The name of the input.
       */
      name: string;
      /**
       * The title of the input.
       */
      title: string;
      /**
       * The type of the input.
       */
      type: "textfield";
      /**
       * The placeholder of the input.
       */
      placeholder?: string;
      /**
       * Whether the input is optional.
       */
      optional?: boolean;
      /**
       * The default value of the input.
       */
      default?: string;
      /**
       * Whether the input should be secure.
       */
      secure?: boolean;
    }
  | {
      /**
       * The name of the input.
       */
      name: string;
      /**
       * The title of the input.
       */
      title: string;
      /**
       * Whether the input is optional.
       */
      optional?: boolean;
      /**
       * The type of the input.
       */
      type: "checkbox";
      /**
       * The default value of the input.
       */
      default?: boolean;
      /**
       * The label of the input.
       */
      label?: string;
      /**
       * The text substitution to use when the input is true.
       */
      trueSubstitution?: string;
      /**
       * The text substitution to use when the input is false.
       */
      falseSubstitution?: string;
    }
  | {
      /**
       * The name of the input.
       */
      name: string;
      /**
       * The title of the input.
       */
      title: string;
      /**
       * The type of the input.
       */
      type: "textarea";
      /**
       * Whether the input is optional.
       */
      optional?: boolean;
      /**
       * The placeholder of the input.
       */
      placeholder?: string;
      /**
       * The default value of the input.
       */
      default?: string;
    }
  | {
      /**
       * The name of the input.
       */
      name: string;
      /**
       * The title of the input.
       */
      title: string;
      /**
       * Whether the input is optional.
       */
      optional?: boolean;
      /**
       * The type of the input.
       */
      type: "dropdown";
      /**
       * The items of the input.
       */
      items: {
        /**
         * The title of the item.
         */
        title: string;
        /**
         * The value of the item.
         */
        value: string;
      }[];
      /**
       * The default value of the input.
       */
      default?: string;
    };
export type OnSuccess = "copy" | "paste" | "open" | "reload";

export interface List {
  /**
   * The type of the response.
   */
  type: "list";
  /**
   * The title of the page.
   */
  title?: string;
  onQueryChange?: PageProvider;
  emptyView?: {
    /**
     * The text to show when the list is empty.
     */
    text: string;
    /**
     * The actions to show when the list is empty.
     */
    actions?: Action[];
  };
  /**
   * Whether to show the detail on the right side of the list.
   */
  showDetail?: boolean;
  items?: Listitem[] | null;
}
export interface PageProvider {
  command?: Command;
  request?: Request;
  text?: string;
  expression?: Expression;
  [k: string]: unknown;
}
export interface Listitem {
  /**
   * The title of the item.
   */
  title: string;
  /**
   * The id of the item.
   */
  id?: string;
  /**
   * The subtitle of the item.
   */
  subtitle?: string;
  detail?: {
    command?: Command;
    request?: Request;
    text?: string;
    expression?: Expression;
  };
  /**
   * The accessories to show on the right side of the item.
   */
  accessories?: string[];
  /**
   * The actions attached to the item.
   */
  actions?: Action[];
}
/**
 * A detail view displaying a preview and actions.
 */
export interface Detail {
  /**
   * The type of the response.
   */
  type: "detail";
  /**
   * The title of the page.
   */
  title?: string;
  command?: Command;
  request?: Request;
  text?: string;
  expression?: Expression;
  /**
   * The actions attached to the detail view.
   */
  actions?: Action[];
}
export interface Form {
  /**
   * The type of the response.
   */
  type: "form";
  /**
   * The title of the page.
   */
  title?: string;
  submitAction: Action;
}



export function List(list: List) {
  return list;
}

export function Detail(detail: Detail) {
  return detail;
}

export function Item(item: Listitem) {
  return item;
}

export function Action(action: Action) {
  return action;
}

export function Input(input: Input) {
  return input;
}
