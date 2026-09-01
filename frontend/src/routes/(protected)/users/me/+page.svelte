<script lang="ts">
	import { page } from "$app/state";
	import { focusTrap } from "$lib/utils/a11y";
	import { slide } from "svelte/transition";
	import { MeController } from "./me_controller.svelte";
	import Button from "$lib/components/ui/Button.svelte";
	import Modal from "$lib/components/ui/Modal.svelte";
	import Input from "$lib/components/ui/Input.svelte";
	import CodeInput from "$lib/components/ui/CodeInput.svelte";
	import ValidationCard from "$lib/components/ui/ValidationCard.svelte";

	const isPasswordSetted = $derived(page.data.user?.hasPassword ?? true);
	const is2FAEnabled = $derived(page.data.user?.is2FAEnabled ?? false);

	const controller = new MeController(
		() => is2FAEnabled,
		() => isPasswordSetted
	);

	type ModalType = "2fa" | "save-password" | null;
	let activeModal = $state<ModalType>(null);
	let showPassword = $state(false);
	let showNewPassword = $state(false);
	let showConfirmPassword = $state(false);

	let touched = $state({
		password: false,
		newPassword: false,
		confirmPassword: false,
		code: false
	});

	const modal2FAVariant = $derived(is2FAEnabled ? "danger" : "success");
	const modal2FATitle = $derived(
		is2FAEnabled
			? "Desativar autenticação em dois fatores?"
			: "Ativar autenticação em dois fatores?"
	);

	const twoFactorLabels: Record<string, string> = $derived({
		generateCode: `Você precisará de um aplicativo autenticador. 
		${isPasswordSetted ? "Insira sua senha e clique" : "Clique"} em gerar chave.`,
		saveCode: "Salve a chave em seu aplicativo autenticador e insira o código gerado nele."
	});

	const modal2FADescription = $derived(
		is2FAEnabled
			? `Desativar a autenticação em dois fatores deixará sua conta vulnerável! Caso deseje continuar
			${isPasswordSetted ? "insira sua senha novamente e " : ""} clique em desativar.`
			: !controller.twoFactorSecret
				? twoFactorLabels.generateCode
				: twoFactorLabels.saveCode
	);

	const modalSavePasswordTitle = $derived(isPasswordSetted ? "Alterar senha" : "Definir senha");
	const modalSavePasswordDescription = $derived(
		isPasswordSetted
			? "Para mudar sua senha digite sua senha antiga e a nova senha. Você será desconectado de outros dispositivos!"
			: "Para definir uma senha para sua conta basta digita-la abaixo. Você será desconectado de outros dispositivos!"
	);

	const openModal = (name: ModalType) => (activeModal = name);
	const closeModal = () => (activeModal = null);

	const togglePassword = () => (showPassword = !showPassword);
	const toggleNewPassword = () => (showNewPassword = !showNewPassword);
	const toggleConfirmPassword = () => (showConfirmPassword = !showConfirmPassword);

	async function handleGenerate2FA(e: SubmitEvent) {
		e.preventDefault();
		const success = await controller.generate2FA();
		if (success) {
			touched.password = false;
		}
	}

	async function handleSave2FA(e: SubmitEvent) {
		e.preventDefault();
		const success = await controller.save2FA();
		if (success) {
			activeModal = null;
			touched.code = false;
		}
	}

	async function handleDisable2FA(e: SubmitEvent) {
		e.preventDefault();
		const success = await controller.disable2FA();
		if (success) {
			activeModal = null;
			touched.password = false;
		}
	}

	async function handleSavePassword(e: SubmitEvent) {
		e.preventDefault();
		const success: Boolean = isPasswordSetted
			? await controller.changePassword()
			: await controller.setPassword();
		if (success) {
			activeModal = null;
			touched.password = false;
			touched.newPassword = false;
			touched.confirmPassword = false;
		}
	}
</script>

<svelte:head>
	<title>Minha conta</title>
	<meta name="robots" content="noindex, nofollow" />
	<meta name="description" content="Informações e configurações da conta." />
</svelte:head>

<div class="max-w-4xl space-y-8 font-inter">
	<header
		class="flex flex-col justify-between gap-4 border-b border-app-border pb-6 md:flex-row md:items-center"
	>
		<div>
			<h1 class="font-montserrat text-2xl font-bold tracking-tight text-text-primary md:text-3xl">
				Minha conta
			</h1>
			<p class="mt-1 text-sm text-text-muted">
				Ajuste suas informações pessoais e configurações de segurança.
			</p>
		</div>
	</header>

	<!-- Account -->
	<section class="space-y-6 rounded-xl border border-app-border bg-app-surface p-6 shadow-xl">
		<div class="border-b border-app-border/80 pb-4">
			<h2 class="font-montserrat text-lg font-semibold text-text-primary">Perfil</h2>
			<p class="mt-0.5 text-xs text-text-muted">Informações visíveis do seu perfil.</p>
		</div>

		<form onsubmit={(e) => e.preventDefault()} class="space-y-5">
			<div class="grid grid-cols-1 gap-5 md:grid-cols-2">
				<!-- Name -->
				<div class="space-y-2">
					<label for="name" class="block text-sm font-medium text-text-secondary"
						>Nome completo</label
					>
					<input
						type="text"
						id="name"
						name="name"
						placeholder="Seu nome"
						value={page.data.user?.name}
						class="w-full rounded-lg border border-app-border bg-app-bg/50 px-4 py-2.5 text-sm text-text-primary placeholder-text-muted transition-all focus:border-text-brand focus:bg-app-bg focus:ring-1 focus:ring-text-brand focus:outline-none"
					/>
				</div>

				<!-- Username -->
				<div class="space-y-2">
					<label for="username" class="block text-sm font-medium text-text-secondary"
						>Nome de usuário</label
					>
					<div class="relative flex items-center">
						<span
							class="pointer-events-none absolute left-3.5 text-sm leading-none font-medium text-text-muted select-none"
							>@</span
						>
						<input
							type="text"
							id="username"
							name="username"
							placeholder="dev_user"
							value={page.data.user?.username}
							class="w-full rounded-lg border border-app-border bg-app-bg/50 py-2.5 pr-4 pl-8 text-sm text-text-primary placeholder-text-muted transition-all focus:border-text-brand focus:bg-app-bg focus:ring-1 focus:ring-text-brand focus:outline-none"
						/>
					</div>
				</div>
			</div>

			<!-- Email -->
			<div class="space-y-2">
				<div class="flex items-center justify-between">
					<label for="email" class="block text-sm font-medium text-text-secondary">E-mail</label>
					<span class="text-xs text-text-muted">O e-mail não pode ser alterado</span>
				</div>
				<div class="relative flex items-center">
					<input
						type="email"
						id="email"
						name="email"
						value={page.data.user?.email}
						disabled
						class="w-full cursor-not-allowed rounded-lg border border-app-border/60 bg-app-bg/30 px-4 py-2.5 text-sm text-text-muted opacity-80 select-none"
					/>
					<svg
						class="absolute right-3.5 h-4 w-4 text-text-muted"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
						/>
					</svg>
				</div>
			</div>

			<div class="flex justify-end pt-2">
				<button
					type="submit"
					class="cursor-pointer rounded-lg bg-text-brand px-5 py-2.5 text-sm font-semibold text-app-bg transition-all hover:opacity-90 active:scale-[0.98]"
				>
					Salvar perfil
				</button>
			</div>
		</form>
	</section>

	<!-- Social accounts -->
	<section class="space-y-6 rounded-xl border border-app-border bg-app-surface p-6 shadow-xl">
		<div class="border-b border-app-border/80 pb-4">
			<h2 class="font-montserrat text-lg font-semibold text-text-primary">Contas vinculadas</h2>
			<p class="mt-0.5 text-xs text-text-muted">Conecte suas redes para facilitar o login.</p>
		</div>

		<div class="space-y-3">
			<!-- Google -->
			<div
				class="flex items-center justify-between rounded-lg border border-app-border bg-app-bg/40 p-4 transition-all hover:border-app-border"
			>
				<div class="flex items-center gap-3">
					<div class="shrink-0 rounded-md border border-app-border bg-app-bg p-2">
						<svg class="h-5 w-5" viewBox="0 0 24 24">
							<path
								fill="#EA4335"
								d="M12 5c1.6 0 3 .6 4.1 1.6l3.1-3.1C17.3 1.7 14.8 1 12 1 7.5 1 3.7 3.6 1.9 7.3l3.7 2.9C6.5 7.2 9 5 12 5z"
							/>
							<path
								fill="#4285F4"
								d="M23.5 12.3c0-.8-.1-1.6-.2-2.3H12v4.5h6.5c-.3 1.5-1.1 2.8-2.4 3.7l3.7 2.9c2.2-2 3.7-5 3.7-8.8z"
							/>
							<path
								fill="#FBBC05"
								d="M5.6 14.8c-.2-.7-.4-1.5-.4-2.3s.2-1.6.4-2.3L1.9 7.3C.7 9.7 0 10.8 0 12s.7 2.3 1.9 4.7l3.7-2.9z"
							/>
							<path
								fill="#34A853"
								d="M12 23c3.2 0 6-1.1 8-3l-3.7-2.9c-1.1.7-2.5 1.2-4.3 1.2-3 0-5.5-2.2-6.4-5.2L1.9 16C3.7 19.7 7.5 23 12 23z"
							/>
						</svg>
					</div>
					<div>
						<p class="text-sm font-medium text-text-primary">Google</p>
						<p class="text-xs text-text-muted">Conectado como dev@exemplo.com</p>
					</div>
				</div>

				<!-- Connected state -->
				<button
					type="button"
					class="cursor-pointer rounded-md border border-red-900/60 px-3 py-1.5 text-xs font-medium text-red-400 transition-all hover:border-red-800 hover:bg-red-950/30"
				>
					Desvincular
				</button>
			</div>

			<!-- GitHub -->
			<div
				class="flex items-center justify-between rounded-lg border border-app-border bg-app-bg/40 p-4 transition-all hover:border-app-border"
			>
				<div class="flex items-center gap-3">
					<div class="shrink-0 rounded-md border border-app-border bg-app-bg p-2">
						<svg class="h-5 w-5 fill-current text-text-primary" viewBox="0 0 24 24">
							<path
								d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"
							/>
						</svg>
					</div>
					<div>
						<p class="text-sm font-medium text-text-primary">GitHub</p>
						<p class="text-xs text-text-muted">Não conectado</p>
					</div>
				</div>

				<!-- Disconnected state -->
				<button
					type="button"
					class="cursor-pointer rounded-md border border-app-border bg-app-bg px-3 py-1.5 text-xs font-medium text-text-secondary transition-all hover:border-text-muted hover:text-text-primary"
				>
					Vincular
				</button>
			</div>
		</div>
	</section>

	<!-- Security -->
	<section class="space-y-6 rounded-xl border border-app-border bg-app-surface p-6 shadow-xl">
		<div class="border-b border-app-border/80 pb-4">
			<h2 class="font-montserrat text-lg font-semibold text-text-primary">Segurança</h2>
			<p class="mt-0.5 text-xs text-text-muted">
				Gerencie sua senha e autenticação de dois fatores.
			</p>
		</div>

		<div class="space-y-6">
			<!-- Change password -->
			<div
				class="flex flex-col justify-between gap-4 rounded-lg border border-app-border bg-app-bg/40 p-4 sm:flex-row sm:items-center"
			>
				<div>
					<p class="text-sm font-medium text-text-primary">Senha de acesso</p>
					<p class="mt-0.5 text-xs text-text-muted">Última alteração há mais de 30 dias</p>
				</div>
				<Button
					variant="outline"
					aria-haspopup="dialog"
					disabled={controller.isLoading}
					aria-expanded={activeModal === "save-password"}
					onclick={() => openModal("save-password")}
				>
					<svg
						class="h-4 w-4 text-text-muted"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 0121 9z"
						/>
					</svg>
					<span>{isPasswordSetted ? "Alterar senha" : "Definir senha"}</span>
				</Button>
			</div>

			<!-- 2FA -->
			<div
				class="flex flex-col justify-between gap-4 rounded-lg border border-app-border bg-app-bg/40 p-4 sm:flex-row sm:items-center"
			>
				<div class="space-y-1">
					<div class="flex items-center gap-2">
						<p class="text-sm font-medium text-text-primary">Autenticação de Dois Fatores (2FA)</p>
						<span
							class="invisible rounded-md border px-2 py-0.5 text-xs font-bold uppercase ml:visible
							{page.data.user?.is2FAEnabled
								? 'border-emerald-800 bg-emerald-950 text-emerald-300'
								: 'border-amber-800 bg-amber-950 text-amber-300'}
							 "
						>
							{page.data.user?.is2FAEnabled ? "Ativado" : "Desativado"}
						</span>
					</div>
					<p class="max-w-md text-xs text-text-muted">
						Adicione uma camada extra de segurança usando um aplicativo autenticador (Google
						Authenticator, Authy, etc).
					</p>
				</div>
				<Button
					variant={page.data.user?.is2FAEnabled ? "danger-soft" : "success-soft"}
					aria-haspopup="dialog"
					disabled={controller.isLoading}
					aria-expanded={activeModal === "2fa"}
					onclick={() => openModal("2fa")}
				>
					<svg
						class="h-4 w-4"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
						/>
					</svg>
					<span>{page.data.user?.is2FAEnabled ? "Desativar 2FA" : "Ativar 2FA"}</span>
				</Button>
			</div>
		</div>
	</section>
</div>

{#if activeModal === "2fa"}
	<Modal
		isOpen={activeModal === "2fa"}
		title={modal2FATitle}
		description={modal2FADescription}
		variant={modal2FAVariant}
		isLoading={controller.isLoading}
		onClose={() => closeModal()}
		{focusTrap}
	>
		{#if !is2FAEnabled}
			{#if !controller.twoFactorSecret}
				<form onsubmit={(e) => handleGenerate2FA(e)} class="space-y-5">
					{#if isPasswordSetted}
						<Input
							id="password"
							name="password"
							type={showPassword ? "text" : "password"}
							label="Senha"
							placeholder="••••••••"
							autocomplete="current-password"
							minlength={8}
							required
							disabled={controller.isLoading}
							bind:value={controller.password}
							touched={touched.password}
							error={(touched.password && !controller.isPasswordValid) ||
							controller.apiError?.code === "AUTH_INCORRECT_PASSWORD"
								? controller.apiError?.code === "AUTH_INCORRECT_PASSWORD"
									? "Senha incorreta."
									: "A senha deve conter no mínimo 8 caracteres."
								: undefined}
							onblur={() => (touched.password = true)}
						>
							{#snippet suffixIcon()}
								<button
									type="button"
									onclick={() => togglePassword()}
									class="rounded p-1 text-text-muted transition-colors hover:text-text-primary focus:ring-1 focus:ring-text-brand focus:outline-none"
									aria-label={showPassword ? "Ocultar senha" : "Mostrar senha"}
								>
									{#if showPassword}
										<svg
											class="h-5 w-5"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
											stroke-linejoin="round"
											aria-hidden="true"
										>
											<path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
											<circle cx="12" cy="12" r="3" />
										</svg>
									{:else}
										<svg
											class="h-5 w-5"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
											stroke-linejoin="round"
											aria-hidden="true"
										>
											<path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
											<path
												d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
											/>
											<path
												d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"
											/>
											<line x1="2" x2="22" y1="2" y2="22" />
										</svg>
									{/if}
								</button>
							{/snippet}
						</Input>
					{/if}
					<div
						class="flex flex-col-reverse justify-end gap-3 border-t border-app-border pt-4 sm:flex-row"
					>
						<Button
							type="button"
							variant="dark"
							size="md"
							disabled={controller.isLoading}
							onclick={() => closeModal()}
							class="w-full sm:w-auto"
						>
							Cancelar
						</Button>
						<Button
							type="submit"
							variant="success-soft"
							size="md"
							disabled={controller.isLoading}
							class="w-full sm:w-auto"
						>
							<span>{controller.isLoading ? "Gerando..." : "Gerar chave"}</span>
						</Button>
					</div>
				</form>
			{:else}
				<form onsubmit={(e) => handleSave2FA(e)} class="space-y-5 font-inter">
					<div
						transition:slide={{ duration: 250 }}
						class="flex flex-col items-center gap-5 rounded-lg border border-app-border bg-app-bg/20 p-4"
					>
						<div class="shrink-0 rounded-xl bg-white p-2.5 shadow-inner">
							<img
								src={controller.qrCodeUrl}
								alt="QR Code para Autenticação 2FA"
								class="h-48 w-48 rounded"
							/>
						</div>

						<div class="w-full space-y-2">
							<label for="twoFactorSecret" class="text-xs font-medium text-text-secondary">
								Chave secreta (se não conseguir escanear):
							</label>
							<div class="flex items-center gap-2">
								<div
									class="flex-1 rounded-lg border border-app-border bg-app-bg p-2.5 shadow-inner"
								>
									<code
										id="twoFactorSecret"
										aria-label="Chave secreta de configuração"
										class="font-mono text-sm leading-relaxed font-medium tracking-wider break-all text-emerald-300"
									>
										{controller.twoFactorSecret || "Gerando chave..."}
									</code>
								</div>

								<button
									type="button"
									title="Copiar chave secreta"
									aria-label="Copiar chave secreta para a área de transferência"
									class="rounded-lg border border-app-border bg-app-bg p-2.5 text-text-muted transition-colors hover:bg-app-surface hover:text-text-primary focus:ring-2 focus:ring-emerald-500 focus:outline-none"
								>
									<svg
										class="h-5 w-5"
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
										aria-hidden="true"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
										/>
									</svg>
								</button>
							</div>
						</div>
					</div>

					<CodeInput
						bind:value={controller.code}
						touched={touched.code}
						error={!controller.isCodeValid || controller.apiError?.code === "2FA_INVALID_CODE"
							? "O código deve conter exatamente 6 números."
							: undefined}
						disabled={controller.isLoading}
						onblur={() => (touched.code = true)}
					/>

					<div
						class="flex flex-col-reverse justify-end gap-3 border-t border-app-border pt-4 sm:flex-row"
					>
						<Button
							type="button"
							variant="dark"
							size="md"
							disabled={controller.isLoading}
							onclick={() => closeModal()}
							class="w-full sm:w-auto"
						>
							Cancelar
						</Button>
						<Button
							type="submit"
							variant="success-soft"
							size="md"
							disabled={controller.isLoading}
							class="w-full sm:w-auto"
						>
							<span>{controller.isLoading ? "Salvando..." : "Salvar chave"}</span>
						</Button>
					</div>
				</form>
			{/if}
		{:else}
			<form onsubmit={(e) => handleDisable2FA(e)} class="space-y-5 font-inter">
				{#if isPasswordSetted}
					<Input
						id="password"
						name="password"
						type={showPassword ? "text" : "password"}
						label="Senha"
						placeholder="••••••••"
						autocomplete="current-password"
						minlength={8}
						required
						disabled={controller.isLoading}
						bind:value={controller.password}
						touched={touched.password}
						error={(touched.password && !controller.isPasswordValid) ||
						controller.apiError?.code === "AUTH_INCORRECT_PASSWORD"
							? controller.apiError?.code === "AUTH_INCORRECT_PASSWORD"
								? "Senha incorreta."
								: "A senha deve conter no mínimo 8 caracteres."
							: undefined}
						onblur={() => (touched.password = true)}
					>
						{#snippet suffixIcon()}
							<button
								type="button"
								onclick={() => togglePassword()}
								class="rounded p-1 text-text-muted transition-colors hover:text-text-primary focus:ring-1 focus:ring-text-brand focus:outline-none"
								aria-label={showPassword ? "Ocultar senha" : "Mostrar senha"}
							>
								{#if showPassword}
									<svg
										class="h-5 w-5"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										aria-hidden="true"
									>
										<path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
										<circle cx="12" cy="12" r="3" />
									</svg>
								{:else}
									<svg
										class="h-5 w-5"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										aria-hidden="true"
									>
										<path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
										<path
											d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
										/>
										<path
											d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"
										/>
										<line x1="2" x2="22" y1="2" y2="22" />
									</svg>
								{/if}
							</button>
						{/snippet}
					</Input>
				{/if}

				<div
					class="flex flex-col-reverse justify-end gap-3 border-t border-app-border pt-4 sm:flex-row"
				>
					<Button
						type="button"
						variant="dark"
						size="md"
						disabled={controller.isLoading}
						onclick={() => closeModal()}
						class="w-full sm:w-auto"
					>
						Cancelar
					</Button>
					<Button
						type="submit"
						variant="danger"
						size="md"
						disabled={controller.isLoading}
						class="w-full sm:w-auto"
					>
						<span>{controller.isLoading ? "Desativando..." : "Desativar"}</span>
					</Button>
				</div>
			</form>
		{/if}
	</Modal>
{/if}

{#if activeModal === "save-password"}
	<Modal
		isOpen={activeModal === "save-password"}
		title={modalSavePasswordTitle}
		description={modalSavePasswordDescription}
		variant="warning"
		isLoading={controller.isLoading}
		onClose={() => closeModal()}
		{focusTrap}
	>
		<form onsubmit={(e) => handleSavePassword(e)} class="space-y-5">
			{#if isPasswordSetted}
				<Input
					id="password"
					name="password"
					type={showPassword ? "text" : "password"}
					label="Senha"
					placeholder="••••••••"
					autocomplete="current-password"
					minlength={8}
					required
					disabled={controller.isLoading}
					bind:value={controller.password}
					touched={touched.password}
					error={(touched.password && !controller.isPasswordValid) ||
					controller.apiError?.code === "AUTH_INCORRECT_PASSWORD"
						? controller.apiError?.code === "AUTH_INCORRECT_PASSWORD"
							? "Senha incorreta."
							: "A senha deve conter no mínimo 8 caracteres."
						: undefined}
					onblur={() => (touched.password = true)}
				>
					{#snippet suffixIcon()}
						<button
							type="button"
							onclick={() => togglePassword()}
							class="rounded p-1 text-text-muted transition-colors hover:text-text-primary focus:ring-1 focus:ring-text-brand focus:outline-none"
							aria-label={showPassword ? "Ocultar senha" : "Mostrar senha"}
						>
							{#if showPassword}
								<svg
									class="h-5 w-5"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									aria-hidden="true"
								>
									<path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
									<circle cx="12" cy="12" r="3" />
								</svg>
							{:else}
								<svg
									class="h-5 w-5"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									aria-hidden="true"
								>
									<path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
									<path
										d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
									/>
									<path d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
									<line x1="2" x2="22" y1="2" y2="22" />
								</svg>
							{/if}
						</button>
					{/snippet}
				</Input>
			{/if}

			<!-- New Password -->
			<Input
				id="newPassword"
				name="newPassword"
				type={showNewPassword ? "text" : "password"}
				label="Nova senha"
				placeholder="••••••••"
				autocomplete="new-password"
				minlength={8}
				required
				disabled={controller.isLoading}
				bind:value={controller.newPassword}
				touched={touched.newPassword}
				error={touched.newPassword && !controller.isNewPasswordValid
					? "A senha não atende aos requisitos mínimos."
					: undefined}
				onblur={() => (touched.newPassword = true)}
			>
				{#snippet suffixIcon()}
					<button
						type="button"
						onclick={() => toggleNewPassword()}
						class="rounded p-1 text-text-muted transition-colors hover:text-text-primary focus:ring-1 focus:ring-text-brand focus:outline-none"
						aria-label={showNewPassword ? "Ocultar senha" : "Mostrar senha"}
					>
						{#if showNewPassword}
							<svg
								class="h-5 w-5"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								aria-hidden="true"
							>
								<path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
								<circle cx="12" cy="12" r="3" />
							</svg>
						{:else}
							<svg
								class="h-5 w-5"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								aria-hidden="true"
							>
								<path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
								<path
									d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
								/>
								<path d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
								<line x1="2" x2="22" y1="2" y2="22" />
							</svg>
						{/if}
					</button>
				{/snippet}
			</Input>

			{#if controller.newPassword.length > 0}
				<ValidationCard
					title="A senha precisa conter:"
					requirements={[
						{ label: "Pelo menos 8 caracteres", met: controller.hasMinLength },
						{
							label: "Letras maiúsculas e minúsculas",
							met: controller.hasUppercase && controller.hasLowercase
						},
						{ label: "Pelo menos um número (0-9)", met: controller.hasNumber },
						{ label: "Pelo menos um símbolo (@$!%*?&)", met: controller.hasSpecialChar }
					]}
				/>
			{/if}

			<!-- Confirm Password -->
			<Input
				id="confirmPassword"
				name="confirmPassword"
				type={showConfirmPassword ? "text" : "password"}
				label="Confirmar senha"
				placeholder="••••••••"
				autocomplete="new-password"
				required
				bind:value={controller.confirmPassword}
				touched={touched.confirmPassword}
				error={controller.apiError?.code === "USER_PASSWORDS_DONT_MATCH"
					? "As senhas não coincidem."
					: touched.confirmPassword &&
						  controller.confirmPassword.length > 0 &&
						  !controller.isPasswordsMatching
						? "As senhas não coincidem."
						: undefined}
				disabled={controller.isLoading}
				onblur={() => (touched.confirmPassword = true)}
			>
				{#snippet suffixIcon()}
					<button
						type="button"
						onclick={() => toggleConfirmPassword()}
						class="rounded p-1 text-text-muted transition-colors hover:text-text-primary focus:ring-1 focus:ring-text-brand focus:outline-none"
						aria-label={showConfirmPassword
							? "Ocultar confirmação de senha"
							: "Mostrar confirmação de senha"}
					>
						{#if showConfirmPassword}
							<svg
								class="h-5 w-5"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								aria-hidden="true"
							>
								<path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
								<circle cx="12" cy="12" r="3" />
							</svg>
						{:else}
							<svg
								class="h-5 w-5"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								aria-hidden="true"
							>
								<path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
								<path
									d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
								/>
								<path d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
								<line x1="2" x2="22" y1="2" y2="22" />
							</svg>
						{/if}
					</button>
				{/snippet}
			</Input>

			<!-- Modal Footer -->
			<div
				class="flex flex-col-reverse justify-end gap-3 border-t border-app-border pt-4 sm:flex-row"
			>
				<Button
					type="button"
					variant="dark"
					size="md"
					disabled={controller.isLoading}
					onclick={() => closeModal()}
					class="w-full sm:w-auto"
				>
					Cancelar
				</Button>
				<Button
					type="submit"
					variant="success-soft"
					size="md"
					disabled={controller.isLoading}
					class="w-full sm:w-auto"
				>
					<span>{controller.isLoading ? "Salvando..." : "Salvar senha"}</span>
				</Button>
			</div>
		</form>
	</Modal>
{/if}
