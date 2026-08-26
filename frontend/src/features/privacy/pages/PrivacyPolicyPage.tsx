import { Link, Paper, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { privacyPolicyPageStyles } from '../styles/privacyPolicyPage.styles';

const policySections = [
  {
    title: 'Dados tratados',
    body:
      'A aplicação trata dados de administradores, como nome, e-mail, senha protegida por hash e dados de login Google quando usado. Também trata estruturas de formulários e respostas enviadas por visitantes.',
  },
  {
    title: 'Finalidades',
    body:
      'Os dados são usados para autenticar administradores, criar e publicar formulários, receber respostas públicas, listar respostas no painel administrativo e permitir exportação ou exclusão quando necessário.',
  },
  {
    title: 'Formulários públicos',
    body:
      'Cada formulário publicado deve informar finalidade do tratamento, retenção das respostas e e-mail de contato do controlador. O visitante precisa confirmar ciência desse aviso antes de enviar uma resposta.',
  },
  {
    title: 'Cookies essenciais',
    body:
      'O sistema usa cookies essenciais de sessão e de segurança do fluxo Google OAuth. Eles mantêm o administrador autenticado e ajudam a validar o retorno do login externo.',
  },
  {
    title: 'Compartilhamento',
    body:
      'Quando o login Google é usado, o backend consulta dados básicos do perfil Google, como identificador, e-mail verificado e nome. A aplicação não inclui rastreadores ou analytics no frontend.',
  },
  {
    title: 'Direitos do titular',
    body:
      'O titular pode solicitar informações, acesso, correção, exportação ou exclusão dos seus dados ao controlador informado no formulário. O painel administrativo oferece exportação e exclusão de respostas para apoiar esse atendimento.',
  },
  {
    title: 'Retenção e segurança',
    body:
      'A retenção das respostas segue o prazo informado em cada formulário. Sessões possuem expiração, senhas usam hash bcrypt e tokens de sessão são armazenados no banco apenas como hash HMAC.',
  },
];

export default function PrivacyPolicyPage() {
  return (
    <Paper sx={privacyPolicyPageStyles.panel}>
      <Stack sx={privacyPolicyPageStyles.stack}>
        <Stack sx={privacyPolicyPageStyles.header}>
          <Typography variant="h4">Política de Privacidade</Typography>
          <Typography color="text.secondary">
            Informações básicas sobre tratamento de dados pessoais no Form Builder.
          </Typography>
        </Stack>

        {policySections.map((section) => (
          <Stack key={section.title} sx={privacyPolicyPageStyles.section}>
            <Typography variant="h6">{section.title}</Typography>
            <Typography color="text.secondary">{section.body}</Typography>
          </Stack>
        ))}

        <Typography variant="body2" color="text.secondary">
          Para formulários específicos, use o contato informado no próprio formulário. Para voltar ao sistema, acesse{' '}
          <Link component={RouterLink} to="/">
            Entrar
          </Link>
          .
        </Typography>
      </Stack>
    </Paper>
  );
}
