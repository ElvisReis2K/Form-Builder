export const guestNavItems = [
  { label: 'Entrar', to: '/' },
  { label: 'Privacidade', to: '/privacidade' },
] as const;

export const authenticatedNavItems = [
  { label: 'Formulários', to: '/admin' },
  { label: 'Administração', to: '/admin/workspace' },
  { label: 'Privacidade', to: '/privacidade' },
] as const;
