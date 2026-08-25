import { Route, Routes } from 'react-router-dom';

import LoginPage from '../features/auth/pages/LoginPage';
import AdminHomePage from '../features/forms/pages/AdminHomePage';
import PublicFormPage from '../features/forms/pages/PublicFormPage';
import FormResponsesPage from '../features/responses/pages/FormResponsesPage';

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/admin" element={<AdminHomePage />} />
      <Route path="/admin/forms/:formId/responses" element={<FormResponsesPage />} />
      <Route path="/f/:slug" element={<PublicFormPage />} />
    </Routes>
  );
}
