import type { AlgorithmEditor } from "$lib/states/editor.svelte";
import { BaseController } from "./base_controller.svelte";

export abstract class BaseEditorController extends BaseController {
	protected _isSuccess = $state(false);
	protected _link = $state("");

	constructor() {
		super();
	}

	get isSuccess() {
		return this._isSuccess;
	}

	get link() {
		return this._link;
	}

	abstract save(editor: AlgorithmEditor): Promise<boolean>;
}
