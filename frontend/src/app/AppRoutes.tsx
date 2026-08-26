import { Route, Routes } from 'react-router-dom';

import RequireAuth from '../features/auth/components/RequireAuth';
import LoginPage from '../features/auth/pages/LoginPage';
import AdminHomePage from '../features/forms/pages/AdminHomePage';
import PublicFormPage from '../features/forms/pages/PublicFormPage';
import SavedFormsPage from '../features/forms/pages/SavedFormsPage';
import PrivacyPolicyPage from '../features/privacy/pages/PrivacyPolicyPage';
import FormResponsesPage from '../features/responses/pages/FormResponsesPage';

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route
        path="/admin"
        element={
          <RequireAuth>
            <SavedFormsPage />
          </RequireAuth>
        }
      />
      <Route
        path="/admin/workspace"
        element={
          <RequireAuth>
            <AdminHomePage />
          </RequireAuth>
        }
      />
      <Route
        path="/admin/forms/:formId/responses"
        element={
          <RequireAuth>
            <FormResponsesPage />
          </RequireAuth>
        }
      />
      <Route path="/f/:slug" element={<PublicFormPage />} />
      <Route path="/privacidade" element={<PrivacyPolicyPage />} />
    </Routes>
  );
}
