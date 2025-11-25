# Test Cases Documentation

**Proyecto Integrador: Mini Red Social**  
**Alumnos:** Octavio Carpineti - Kevin Massholder  
**Materia:** Ingeniería de Software III  
**TP7 + TP8 Integration**

---

## Tabla de Contenidos

1. [Backend Unit Tests - Auth Service](#backend-unit-tests---auth-service)
2. [Backend Unit Tests - Post Service](#backend-unit-tests---post-service)
3. [Backend Unit Tests - Auth Handler](#backend-unit-tests---auth-handler)
4. [Backend Unit Tests - Post Handler](#backend-unit-tests---post-handler)
5. [Backend Integration Tests](#backend-integration-tests)
6. [Frontend Unit Tests - Components](#frontend-unit-tests---components)
7. [Frontend Unit Tests - Services](#frontend-unit-tests---services)
8. [E2E Tests - Cypress](#e2e-tests---cypress)

---

## Backend Unit Tests - Auth Service

### TestRegister_Success
**Objetivo:** Verificar registro exitoso de usuario nuevo.  
**Escenario:** Email no existe, datos válidos.  
**Resultado esperado:** Usuario creado y retornado.

### TestRegister_EmailVacio
**Objetivo:** Validar que falle cuando email está vacío.  
**Escenario:** Email vacío en request.  
**Resultado esperado:** Error "el email es requerido".

### TestRegister_EmailInvalido
**Objetivo:** Validar formato de email sin @.  
**Escenario:** Email sin arroba.  
**Resultado esperado:** Error "el email debe ser válido".

### TestRegister_PasswordCorto
**Objetivo:** Validar contraseña mínima de 6 caracteres.  
**Escenario:** Password de 3 caracteres.  
**Resultado esperado:** Error longitud contraseña.

### TestRegister_UsernameVacio
**Objetivo:** Validar que username no esté vacío.  
**Escenario:** Username vacío.  
**Resultado esperado:** Error "el nombre de usuario es requerido".

### TestRegister_EmailDuplicado
**Objetivo:** Prevenir email duplicado.  
**Escenario:** Email ya registrado.  
**Resultado esperado:** Error "el email ya está registrado".

### TestLogin_Success
**Objetivo:** Verificar login exitoso.  
**Escenario:** Credenciales existentes.  
**Resultado esperado:** Usuario autenticado.

### TestLogin_EmailVacio
**Objetivo:** Validar email requerido en login.  
**Escenario:** Email vacío.  
**Resultado esperado:** Error "el email es requerido".

### TestLogin_PasswordVacio
**Objetivo:** Validar password requerido.  
**Escenario:** Password vacío.  
**Resultado esperado:** Error "la contraseña es requerida".

### TestLogin_UsuarioNoExiste
**Objetivo:** Manejar usuario inexistente.  
**Escenario:** Email no encontrado.  
**Resultado esperado:** Error "credenciales inválidas".

### TestLogin_PasswordIncorrecta
**Objetivo:** Validar password incorrecta.  
**Escenario:** Password erróneo.  
**Resultado esperado:** Error "credenciales inválidas".

---

## Backend Unit Tests - Post Service

### TestCreatePost_Success
**Objetivo:** Creación exitosa de post.  
**Escenario:** Título y contenido válidos, usuario existe.  
**Resultado esperado:** Post creado con datos correctos.

### TestCreatePost_UserNotFound
**Objetivo:** Validar usuario existente.  
**Escenario:** UserID no existe.  
**Resultado esperado:** Error usuario no encontrado.

### TestCreatePost_RepoError
**Objetivo:** Manejar error de repositorio.  
**Escenario:** Repositorio falla al guardar.  
**Resultado esperado:** Error propagado.

### TestCreatePost_TitleVacio
**Objetivo:** Validar título requerido.  
**Escenario:** Título vacío.  
**Resultado esperado:** Error validación título.

### TestCreatePost_ContentVacio
**Objetivo:** Validar contenido requerido.  
**Escenario:** Contenido vacío.  
**Resultado esperado:** Error validación contenido.

### TestDeletePost_Success
**Objetivo:** Eliminación exitosa por autor.  
**Escenario:** Usuario elimina su post.  
**Resultado esperado:** Post eliminado.

### TestDeletePost_PostNoExiste
**Objetivo:** Manejar post inexistente.  
**Escenario:** ID de post no existe.  
**Resultado esperado:** Error post no encontrado.

### TestDeletePost_NoEsAutor
**Objetivo:** Validar permisos de eliminación.  
**Escenario:** Usuario intenta eliminar post ajeno.  
**Resultado esperado:** Forbidden - no autorizado.

### TestDeleteComment_Success
**Objetivo:** Eliminación exitosa de comentario por autor.  
**Escenario:** Autor elimina su comentario.  
**Resultado esperado:** Comentario eliminado.

### TestDeleteComment_PostNoExiste
**Objetivo:** Validar post existente para comentario.  
**Escenario:** Post del comentario no existe.  
**Resultado esperado:** Error post no encontrado.

### TestDeleteComment_UsuarioNoExiste
**Objetivo:** Validar usuario existente.  
**Escenario:** Usuario del comentario no existe.  
**Resultado esperado:** Error usuario no encontrado.

### TestDeleteComment_NoEsAutor
**Objetivo:** Validar permisos en comentarios.  
**Escenario:** Usuario elimina comentario ajeno.  
**Resultado esperado:** Forbidden - no autorizado.

### TestGetAllPosts_Success
**Objetivo:** Obtener todos los posts.  
**Escenario:** Posts existentes.  
**Resultado esperado:** Lista de posts retornada.

### TestGetAllPosts_Empty
**Objetivo:** Manejar lista vacía de posts.  
**Escenario:** Sin posts.  
**Resultado esperado:** Lista vacía retornada.

### TestGetPostByID_Success
**Objetivo:** Obtener post específico.  
**Escenario:** Post existe.  
**Resultado esperado:** Post retornado.

### TestGetPostByID_InvalidID
**Objetivo:** Validar ID válido para post.  
**Escenario:** ID negativo o cero.  
**Resultado esperado:** Error ID inválido.

### TestGetPostByID_NotFound
**Objetivo:** Manejar post no encontrado.  
**Escenario:** Post no existe.  
**Resultado esperado:** Error post no encontrado.

### TestCreateComment_Success
**Objetivo:** Creación exitosa de comentario.  
**Escenario:** Contenido válido, post y usuario existen.  
**Resultado esperado:** Comentario creado.

### TestCreateComment_EmptyContent
**Objetivo:** Validar contenido de comentario.  
**Escenario:** Contenido vacío.  
**Resultado esperado:** Error validación.

### TestCreateComment_PostNotFound
**Objetivo:** Validar post existente.  
**Escenario:** Post no existe.  
**Resultado esperado:** Error post no encontrado.

### TestCreateComment_UserNotFound
**Objetivo:** Validar usuario existente.  
**Escenario:** Usuario no existe.  
**Resultado esperado:** Error usuario no encontrado.

### TestGetCommentsByPostID_Success
**Objetivo:** Obtener comentarios de un post.  
**Escenario:** Post con comentarios.  
**Resultado esperado:** Lista de comentarios.

### TestGetCommentsByPostID_PostNotFound
**Objetivo:** Manejar post inexistente.  
**Escenario:** Post no existe.  
**Resultado esperado:** Error post no encontrado.

### TestGetCommentsByPostID_Empty
**Objetivo:** Manejar post sin comentarios.  
**Escenario:** Post sin comentarios.  
**Resultado esperado:** Lista vacía.

---

## Backend Unit Tests - Auth Handler

### TestAuthHandler_Register_Success
**Objetivo:** Handler registra usuario exitosamente.  
**Escenario:** JSON válido, servicio registra correctamente.  
**Resultado esperado:** Status 201, usuario retornado.

### TestAuthHandler_Register_InvalidJSON
**Objetivo:** Manejar JSON inválido en registro.  
**Escenario:** Request con JSON malformado.  
**Resultado esperado:** Status 400, error JSON inválido.

### TestAuthHandler_Register_ServiceError
**Objetivo:** Manejar error del servicio.  
**Escenario:** Servicio lanza error.  
**Resultado esperado:** Status 400, error del servicio.

### TestAuthHandler_Login_Success
**Objetivo:** Handler hace login exitosamente.  
**Escenario:** Credenciales válidas.  
**Resultado esperado:** Status 200, usuario retornado.

### TestAuthHandler_Login_InvalidJSON
**Objetivo:** Manejar JSON inválido en login.  
**Escenario:** JSON malformado.  
**Resultado esperado:** Status 400, error JSON inválido.

### TestAuthHandler_Login_ServiceError
**Objetivo:** Manejar error de servicio en login.  
**Escenario:** Servicio falla.  
**Resultado esperado:** Status 401, error del servicio.

---

## Backend Unit Tests - Post Handler

### TestPostHandler_CreatePost_Success
**Objetivo:** Handler crea post exitosamente.  
**Escenario:** JSON válido, headers correctos.  
**Resultado esperado:** Status 201, post creado.

### TestPostHandler_CreatePost_InvalidJSON
**Objetivo:** Validar JSON en creación de post.  
**Escenario:** JSON inválido.  
**Resultado esperado:** Status 400, error JSON.

### TestPostHandler_CreatePost_MissingUserID
**Objetivo:** Validar autenticación.  
**Escenario:** Sin header X-User-ID.  
**Resultado esperado:** Status 401, no autenticado.

### TestPostHandler_CreatePost_InvalidUserID
**Objetivo:** Validar userID como entero.  
**Escenario:** X-User-ID no válido.  
**Resultado esperado:** Status 400, userID inválido.

### TestPostHandler_GetAllPosts_Success
**Objetivo:** Obtener lista de posts via HTTP.  
**Escenario:** Posts existen.  
**Resultado esperado:** Status 200, lista de posts.

### TestPostHandler_GetAllPosts_ServiceError
**Objetivo:** Manejar error de servicio.  
**Escenario:** Servicio falla.  
**Resultado esperado:** Status 500, error interno.

### TestPostHandler_GetPostByID_Success
**Objetivo:** Obtener post específico via URL param.  
**Escenario:** ID válido, post existe.  
**Resultado esperado:** Status 200, post retornado.

### TestPostHandler_GetPostByID_InvalidID
**Objetivo:** Validar ID en URL.  
**Escenario:** ID no entero.  
**Resultado esperado:** Status 400, ID inválido.

### TestPostHandler_DeletePost_Success
**Objetivo:** Eliminar post via HTTP DELETE.  
**Escenario:** Usuario autorizado.  
**Resultado esperado:** Status 200, post eliminado.

### TestPostHandler_DeletePost_PermissionDenied
**Objetivo:** Validar permisos en eliminación.  
**Escenario:** Usuario no autor.  
**Resultado esperado:** Status 403, forbidden.

### TestPostHandler_CreateComment_Success
**Objetivo:** Crear comentario en post específico.  
**Scenarios:** JSON válido, post existe.  
**Resultado esperado:** Status 201, comentario creado.

### TestPostHandler_CreateComment_InvalidJSON
**Objetivo:** Validar JSON en comentario.  
**Escenario:** JSON inválido.  
**Resultado esperado:** Status 400, error JSON.

### TestPostHandler_CreateComment_InvalidPostID
**Objetivo:** Validar postID en URL.  
**Escenario:** PostID inválido.  
**Resultado esperado:** Status 400, ID inválido.

### TestPostHandler_GetComments_Success
**Objetivo:** Obtener comentarios de post via HTTP.  
**Escenario:** Post con comentarios.  
**Resultado esperado:** Status 200, lista de comentarios.

### TestPostHandler_GetComments_InvalidPostID
**Objetivo:** Validar postID en URL para comentarios.  
**Escenario:** ID inválido.  
**Resultado esperado:** Status 400, ID inválido.

### TestPostHandler_DeleteComment_Success
**Objetivo:** Eliminar comentario via HTTP.  
**Escenario:** Usuario autorizado.  
**Resultado esperado:** Status 200, comentario eliminado.

### TestPostHandler_DeleteComment_InvalidPostID
**Objetivo:** Validar postID en eliminación comentario.  
**Escenario:** PostID inválido.  
**Resultado esperado:** Status 400, ID inválido.


---

## Frontend Unit Tests - Components

### Login Component Tests
**renderiza el formulario de login correctamente**  
**Objetivo:** Verificar estructura del formulario.  
**Escenario:** Login inicial.  
**Resultado esperado:** Heading, inputs, botón visibles.

**muestra formulario de registro al cambiar modo**  
**Objetivo:** Toggle entre login y registro.  
**Escenario:** Click en link de registro.  
**Resultado esperado:** Mostrar campos de registro.

**login exitoso llama a onLoginSuccess**  
**Objetivo:** Flujo exitoso de login.  
**Escenario:** Mock API success.  
**Resultado esperado:** Llama callback con usuario.

**muestra error cuando login falla**  
**Objetivo:** Mostrar errores de API.  
**Escenario:** Mock API error.  
**Resultado esperado:** Muestra mensaje de error.

**deshabilita el botón mientras está cargando**  
**Objetivo:** UX: prevenir múltiples submissions.  
**Escenario:** Durante loading.  
**Resultado esperado:** Botón disabled.

### CreatePost Component Tests (basado en estructura similar)
- Renderizar formulario correctamente
- Crear post exitosamente
- Mostrar error en validación
- Manejar estado de loading

### PostList Component Tests
- Renderizar lista vacía
- Mostrar posts correctamente
- Manejar estado de carga

### CommentForm & CommentList Tests
- Comentarios funcionales como posts

### PostDetail Tests
- Mostrar post individual
- Integrar comentarios

---

## Frontend Unit Tests - Services

### AuthService Tests (39 tests total - por component)
**login exitoso retorna usuario**  
**Objetivo:** Servicio llama API correcta.  
**Escenario:** Credenciales válidas.  
**Resultado esperado:** Usuario retornado.

**login fallido lanza error**  
**Objetivo:** Manejar errores de API.  
**Scenarios:** 401, 500, etc.  
**Resultado esperado:** Error apropiado.

**registro exitoso**  
**Objetivo:** Registrar nuevo usuario.  
**Escenario:** Datos válidos.  
**Resultado esperado:** Usuario creado.

**registro con email duplicado**  
**Objetivo:** Manejar error de registro.  
**Escenario:** Email ya existe.  
**Resultado esperado:** Error específico.

### PostService Tests
- Obtener posts
- Crear post
- Eliminar post
- Crear comentario
- Obtener comentarios
- Eliminar comentario

---

## E2E Tests - Cypress

### Auth Flow (5 tests)
**debería mostrar el formulario de login por defecto**  
**Objetivo:** Pantalla inicial de login.  
**Escenario:** Visitar /  
**Resultado esperado:** Formulario visible.

**debería cambiar entre login y registro**  
**Objetivo:** Navegación entre modos.  
**Escenario:** Click en links.  
**Resultado esperado:** Cambio de formularios.

**debería mostrar error con credenciales inválidas**  
**Objetivo:** Validar errores de login.  
**Escenario:** Mock API 401.  
**Resultado esperado:** Mensaje de error mostrado.

**debería hacer login exitoso**  
**Objetivo:** Flujo completo de autenticación.  
**Scenario:** Login válido.  
**Resultado esperado:** Redirección a app, usuario logged in.

**debería registrarse exitosamente**  
**Objetivo:** Registro de nuevos usuarios.  
**Escenario:** Formulario registro válido.  
**Resultado esperado:** Usuario creado, logged in.

### Posts Flow (5 tests)
**debería mostrar lista de posts al inicio**  
**Objetivo:** Home page con posts.  
**Escenario:** Usuario logueado.  
**Resultado esperado:** Lista de posts visible.

**debería crear un post correctamente**  
**Objetivo:** Añadir nuevo contenido.  
**Scenarios:** Formulario válido.  
**Resultado esperado:** Post aparece en lista.

**debería validar campos obligatorios en post**  
**Objetivo:** Validación client-side.  
**Scenarios:** Campos vacíos.  
**Resultado esperado:** Errores mostrados.

**debería mostrar error si falla creación de post**  
**Objetivo:** Manejar errores del servidor.  
**Scenarios:** Mock API error.  
**Resultado esperado:** Mensaje de error.

**debería eliminar post por click**  
**Objetivo:** Funcionalidad de eliminación.  
**Scenarios:** Click en eliminar propio post.  
**Resultado esperado:** Post removido de lista.

### Comments Flow (4 tests)
**debería mostrar comentarios en post detal**  
**Objetivo:** Interacción post-comentarios.  
**Scenarios:** Abrir post con comentarios.  
**Resultado esperado:** Lista visible.

**debería crear comentario exitosamente**  
**Objetivo:** Añadir respuesta a post.  
**Scenarios:** Formulario válido en detalle.  
**Resultado esperado:** Comentario añadido.

**debería validar contenido de comentario**  
**Objetivo:** Validación en comentarios.  
**Scenarios:** Contenido vacío.  
**Resultado esperado:** Error mostrado.

**debería eliminar comentario correctamente**  
**Objetivo:** Moderación de comentarios.  
**Scenarios:** Eliminar propio comentario.  
**Resultado esperado:** Comentario removido.

### Full Flow (1 test)
**debería completar flujo completo de usuario**  
**Objetivo:** End-to-end user journey.  
**Scenarios:** Registro → Login → Crear post → Comentar → Eliminar.  
**Resultado esperado:** Todas features funcionan en secuencia.

---

## Cobertura de Tests

- **Backend Services:** 86.5% (validations + business logic)  
- **Backend Handlers:** 50.4% (HTTP responses + middleware)  
- **Repositories Integration:** ~85% (real database operations)  
- **Frontend Components:** 92.44% (UI + state management)  
- **E2E Cypress:** 15 tests (complete user flows)  
- **Total Tests:** 89 unitários + integration + E2E

## Métricas de Calidad Alcanzadas

- ❌ **Coverage thresholds:** Above 70% for all layers  
- ❌ **Integration tests:** Real DB with proper cleanup  
- ❌ **E2E coverage:** All major user stories tested  
- ❌ **Error handling:** Proper HTTP status + JSON errors  
- ❌ **Security:** Authentication & authorization tested
