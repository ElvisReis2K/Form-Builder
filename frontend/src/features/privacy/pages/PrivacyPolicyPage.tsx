import { Link, Paper, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { privacyPolicyPageStyles } from '../styles/privacyPolicyPage.styles';

const policySections = [
  {
    title: 'Dados tratados',
    body:
      'A aplicacao trata dados de administradores, como nome, e-mail, senha protegida por hash e dados de login Google quando usado. Tambem trata estruturas de formularios e respostas enviadas por visitantes.',
  },
  {
    title: 'Finalidades',
    body:
      'Os dados sao usados para autenticar administradores, criar e publicar formularios, receber respostas publicas, listar respostas no painel administrativo e permitir exportacao ou exclusao quando necessario.',
  },
  {
    title: 'Formularios publicos',
    body:
      'Cada formulario publicado deve informar finalidade do tratamento, retencao das respostas e e-mail de contato do controlador. O visitante precisa confirmar ciencia desse aviso antes de enviar uma resposta.',
  },
  {
    title: 'Cookies essenciais',
    body:
      'O sistema usa cookies essenciais de sessao e de seguranca do fluxo Google OAuth. Eles mantem o administrador autenticado e ajudam a validar o retorno do login externo.',
  },
  {
    title: 'Compartilhamento',
    body:
      'Quando o login Google e usado, o backend consulta dados basicos do perfil Google, como identificador, e-mail verificado e nome. A aplicacao nao inclui rastreadores ou analytics no frontend.',
  },
  {
    title: 'Direitos do titular',
    body:
      'O titular pode solicitar informacoes, acesso, correcao, exportacao ou exclusao dos seus dados ao controlador informado no formulario. O painel administrativo oferece exportacao e exclusao de respostas para apoiar esse atendimento.',
  },
  {
    title: 'Retencao e seguranca',
    body:
      'A retencao das respostas segue o prazo informado em cada formulario. Sessoes possuem expiracao, senhas usam hash bcrypt e tokens de sessao sao armazenados no banco apenas como hash HMAC.',
  },
];

export default function PrivacyPolicyPage() {
  return (
    <Paper sx={privacyPolicyPageStyles.panel}>
      <Stack sx={privacyPolicyPageStyles.stack}>
        <Stack sx={privacyPolicyPageStyles.header}>
          <Typography variant="h4">Politica de Privacidade</Typography>
          <Typography color="text.secondary">
            Informacoes basicas sobre tratamento de dados pessoais no Form Builder.
          </Typography>
        </Stack>

        {policySections.map((section) => (
          <Stack key={section.title} sx={privacyPolicyPageStyles.section}>
            <Typography variant="h6">{section.title}</Typography>
            <Typography color="text.secondary">{section.body}</Typography>
          </Stack>
        ))}

        <Typography variant="body2" color="text.secondary">
          Para formularios especificos, use o contato informado no proprio formulario. Para voltar ao sistema, acesse{' '}
          <Link component={RouterLink} to="/">
            Entrar
          </Link>
          .
        </Typography>
      </Stack>
    </Paper>
  );
}
