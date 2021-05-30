<script>
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'

    import { _, locale, setupI18n, isLocaleLoaded } from './i18n'
    import { appRoutes } from './config'
    import Notifications from './components/Notifications.svelte'
    import UserIcon from './components/icons/UserIcon.svelte'
    import GithubIcon from './components/icons/Github.svelte'
    import HollowButton from './components/HollowButton.svelte'
    import LocaleSwitcher from './components/LocaleSwitcher.svelte'
    import GlobalAlerts from './components/GlobalAlerts.svelte'
    import Landing from './pages/Landing.svelte'
    import Retrospectives from './pages/Retrospectives.svelte'
    import Retrospective from './pages/Retrospective.svelte'
    import Organizations from './pages/Organizations.svelte'
    import Organization from './pages/Organization.svelte'
    import Department from './pages/Department.svelte'
    import Team from './pages/Team.svelte'
    import Register from './pages/Register.svelte'
    import Login from './pages/Login.svelte'
    import ResetPassword from './pages/ResetPassword.svelte'
    import UserProfile from './pages/UserProfile.svelte'
    import VerifyAccount from './pages/VerifyAccount.svelte'
    import Admin from './pages/admin/Admin.svelte'
    import AdminUsers from './pages/admin/Users.svelte'
    import AdminOrganizations from './pages/admin/Organizations.svelte'
    import AdminTeams from './pages/admin/Teams.svelte'
    import AdminApikeys from './pages/admin/ApiKeys.svelte'
    import AdminAlerts from './pages/admin/Alerts.svelte'
    import { user } from './stores.js'
    import eventTag from './eventTag.js'
    import apiclient from './apiclient.js'

    setupI18n()

    const { AllowRegistration, AppVersion, PathPrefix } = appConfig
    const footerLinkClasses = 'no-underline text-blue-500 hover:text-blue-800'

    let notifications

    let activeUser
    user.subscribe(w => {
        activeUser = w
    })

    let currentPage = {
        route: Landing,
        params: {},
    }

    const router = Navaid('/')
        .on(appRoutes.landing, () => {
            currentPage = {
                route: Landing,
                params: {},
            }
        })
        .on(`${appRoutes.register}/:retrospectiveId?`, params => {
            currentPage = {
                route: Register,
                params,
            }
        })
        .on(`${appRoutes.login}/:retrospectiveId?`, params => {
            currentPage = {
                route: Login,
                params,
            }
        })
        .on(`${appRoutes.resetPwd}/:resetId`, params => {
            currentPage = {
                route: ResetPassword,
                params,
            }
        })
        .on(`${appRoutes.verifyAct}/:verifyId`, params => {
            currentPage = {
                route: VerifyAccount,
                params,
            }
        })
        .on(appRoutes.profile, params => {
            currentPage = {
                route: UserProfile,
                params,
            }
        })
        .on(appRoutes.retrospectives, () => {
            currentPage = {
                route: Retrospectives,
                params: {},
            }
        })
        .on(`${appRoutes.retrospective}/:retrospectiveId`, params => {
            currentPage = {
                route: Retrospective,
                params,
            }
        })
        .on(appRoutes.organizations, () => {
            currentPage = {
                route: Organizations,
                params: {},
            }
        })
        .on(`${appRoutes.organization}/:organizationId`, params => {
            currentPage = {
                route: Organization,
                params,
            }
        })
        .on(
            `${appRoutes.organization}/:organizationId/team/:teamId`,
            params => {
                currentPage = {
                    route: Team,
                    params,
                }
            },
        )
        .on(
            `${appRoutes.organization}/:organizationId/department/:departmentId`,
            params => {
                currentPage = {
                    route: Department,
                    params,
                }
            },
        )
        .on(
            `${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId`,
            params => {
                currentPage = {
                    route: Team,
                    params,
                }
            },
        )
        .on(`${appRoutes.team}/:teamId`, params => {
            currentPage = {
                route: Team,
                params,
            }
        })
        .on(appRoutes.admin, () => {
            currentPage = {
                route: Admin,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/users`, () => {
            currentPage = {
                route: AdminUsers,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/organizations`, () => {
            currentPage = {
                route: AdminOrganizations,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/teams`, () => {
            currentPage = {
                route: AdminTeams,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/apikeys`, () => {
            currentPage = {
                route: AdminApikeys,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/alerts`, () => {
            currentPage = {
                route: AdminAlerts,
                params: {},
            }
        })
        .listen()

    const xfetch = apiclient(handle401)

    function handle401() {
        eventTag('session_expired', 'engagement', 'unauthorized', () => {
            user.delete()
            router.route(appRoutes.login)
        })
    }

    function logoutUser() {
        xfetch('/api/auth/logout', { method: 'POST' })
            .then(function() {
                eventTag('logout', 'engagement', 'success', () => {
                    user.delete()
                    router.route(appRoutes.landing, true)
                })
            })
            .catch(function(error) {
                notifications.danger(
                    'Error encountered attempting to logout user',
                )
                eventTag('logout', 'engagement', 'failure')
            })
    }

    onDestroy(router.unlisten)
</script>

<style>
    :global(.nav-logo) {
        max-height: 3.75rem;
    }
    :global(.text-brand-color) {
        color: #1c90abff;
    }
    :global(.border-brand-color) {
        border-color: #1c90abff;
    }
    :global(.bg-brand-color) {
        background-color: #1c90abff;
    }
</style>

<Notifications bind:this="{notifications}" />

{#if $isLocaleLoaded}
<div class="flex flex-col min-h-screen">
    <GlobalAlerts registered="{!!activeUser.name}" />

    <nav
        class="flex items-center justify-between flex-wrap px-6 py-2 bg-white"
        role="navigation"
        aria-label="main navigation">
        <div class="flex items-center flex-shrink-0 mr-6">
            <a href="{appRoutes.landing}">
                <img
                    src="{PathPrefix}/img/wakita-logo.png"
                    alt="Wakita"
                    class="nav-logo" />
            </a>
        </div>
        <div class="text-right mt-4 md:mt-0">
            {#if activeUser.name}
                <span class="font-bold mr-2 text-xl">
                    <UserIcon />
                    <a href="{appRoutes.profile}">{activeUser.name}</a>
                </span>
                <HollowButton
                    color="teal"
                    href="{appRoutes.retrospectives}"
                    class="mr-2">
                    My Retro's
                </HollowButton>
                {#if activeUser.type !== 'GUEST'}
                    <HollowButton
                        color="blue"
                        href="{appRoutes.organizations}"
                        class="mr-2">
                        Organizations &amp; Teams
                    </HollowButton>
                {/if}
                {#if !activeUser.type || activeUser.type === 'GUEST'}
                    {#if AllowRegistration}
                        <HollowButton
                            color="orange"
                            href="{appRoutes.register}"
                            class="mr-2">
                            Create Account
                        </HollowButton>
                    {/if}
                    <HollowButton href="{appRoutes.login}">Login</HollowButton>
                {:else}
                    {#if activeUser.type === 'ADMIN'}
                        <HollowButton
                            color="purple"
                            href="{appRoutes.admin}"
                            class="mr-2">
                            Admin
                        </HollowButton>
                    {/if}
                    <HollowButton color="red" onClick="{logoutUser}">
                        Logout
                    </HollowButton>
                {/if}
            {:else}
                {#if AllowRegistration}
                    <HollowButton
                        color="orange"
                        href="{appRoutes.register}"
                        class="mr-2">
                        Create Account
                    </HollowButton>
                {/if}
                <HollowButton href="{appRoutes.login}">Login</HollowButton>
            {/if}
            <LocaleSwitcher
                selectedLocale="{$locale}"
                on:locale-changed="{e => setupI18n({
                        withLocale: e.detail,
                    })}" />
        </div>
    </nav>

    <svelte:component
        this="{currentPage.route}"
        {...currentPage.params}
        {notifications}
        {router}
        {eventTag}
        {xfetch} />

    <footer class="p-6 text-center">
        <GithubIcon />
        <a
            href="https://github.com/StevenWeathers/wakita-retro-tool"
            class="{footerLinkClasses}">
            {$_('appName')}
        </a>
        {@html $_('footer.authoredBy', {
            values: {
                authorOpen: `<a href="http://stevenweathers.com" class="${footerLinkClasses}">`,
                authorClose: `</a>`,
            },
        })}
        {@html $_('footer.license', {
            values: {
                licenseOpen: `<a href="http://www.apache.org/licenses/" class="${footerLinkClasses}">`,
                licenseClose: `</a>`,
            },
        })}
        <br />
        {@html $_('footer.poweredBy', {
            values: {
                svelteOpen: `<a href="https://svelte.dev/" class="${footerLinkClasses}">`,
                svelteClose: `</a>`,
                goOpen: `<a href="https://golang.org/" class="${footerLinkClasses}">`,
                goClose: `</a>`,
            },
        })}
        <div class="text-sm text-gray-500">
            {$_('appVersion', { values: { version: AppVersion } })}
        </div>
    </footer>
</div>
{/if}
