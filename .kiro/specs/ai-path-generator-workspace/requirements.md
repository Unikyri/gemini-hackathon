# Requirements Document

## Introduction

Este documento define los requisitos para dos funcionalidades clave de **Gemini Coding Path**: la generación de rutas de aprendizaje mediante IA (Épica E02) y el workspace de código para resolver ejercicios (Épica E04). Estas funcionalidades permiten que un estudiante ingrese un tema, reciba una ruta personalizada generada por Gemini 3, y pueda trabajar en los ejercicios dentro de un editor integrado.

## Glossary

- **Learning_Path**: Ruta de aprendizaje personalizada generada por IA, compuesta por múltiples nodos de ejercicios.
- **Node**: Unidad individual de aprendizaje dentro de un Path, contiene un ejercicio con enunciado, código boilerplate y tests.
- **Backend_API**: Servidor Go que implementa Clean Architecture, responsable de la lógica de negocio y comunicación con Vertex AI.
- **Vertex_AI**: Servicio de Google Cloud que ejecuta el modelo Gemini 3 para generación de contenido estructurado.
- **System_Prompt**: Instrucciones predefinidas enviadas a Gemini para controlar el formato y estructura de la respuesta.
- **Workspace**: Interfaz de usuario dividida en paneles para resolver ejercicios de código.
- **Monaco_Editor**: Componente de editor de código basado en VS Code con resaltado de sintaxis.
- **Boilerplate_Code**: Código inicial proporcionado al estudiante como punto de partida para un ejercicio.

## Requirements

### Requirement 1: Generación Básica del Path con IA

**User Story:** As a estudiante, I want ingresar un tema en un campo de texto, so that el sistema me genere una ruta de aprendizaje personalizada inicial.

#### Acceptance Criteria

1. WHEN el estudiante envía un prompt con un tema, THE Backend_API SHALL recibir la solicitud y validar que el prompt no esté vacío.
2. WHEN el Backend_API recibe un prompt válido, THE Backend_API SHALL llamar a Vertex_AI con un System_Prompt diseñado para retornar JSON estructurado.
3. WHEN Vertex_AI procesa la solicitud, THE Vertex_AI SHALL retornar una lista de exactamente 5 nodos iniciales.
4. THE Node SHALL contener al menos: id (UUID), título, descripción corta y estado (locked/unlocked).
5. WHEN la generación es exitosa, THE Backend_API SHALL retornar un JSON limpio con path_id, title y la lista de nodes.
6. IF el prompt está vacío o es inválido, THEN THE Backend_API SHALL retornar un error 400 con mensaje descriptivo.
7. IF Vertex_AI falla o no responde, THEN THE Backend_API SHALL retornar un error 503 con mensaje de servicio no disponible.

### Requirement 2: Persistencia del Path Generado

**User Story:** As a sistema, I want guardar la estructura generada en la base de datos, so that no sea necesario volver a llamar a la IA cada vez que el usuario recarga la página.

#### Acceptance Criteria

1. THE Backend_API SHALL crear las tablas learning_paths y path_nodes en PostgreSQL siguiendo el esquema definido en la wiki.
2. WHEN un Path es generado exitosamente, THE Backend_API SHALL insertar el registro en learning_paths con user_id, topic, title y ai_metadata.
3. WHEN un Path es generado exitosamente, THE Backend_API SHALL insertar todos los nodes en path_nodes con path_id, position, title, description_md y status.
4. WHEN un cliente solicita GET /paths/{id}, THE Backend_API SHALL recuperar el Path completo con todos sus nodes desde la base de datos.
5. IF el path_id no existe, THEN THE Backend_API SHALL retornar un error 404 con mensaje descriptivo.

### Requirement 3: Interfaz del Editor (Layout del Workspace)

**User Story:** As a estudiante, I want ver una pantalla dividida con el enunciado y un editor de código, so that tenga todo lo necesario para resolver el ejercicio sin cambiar de ventana.

#### Acceptance Criteria

1. THE Workspace SHALL mostrar un layout de 3 secciones: Instrucciones (izquierda), Editor (centro), Info/Salida (derecha).
2. THE Monaco_Editor SHALL tener resaltado de sintaxis configurado para el lenguaje Go.
3. THE Workspace SHALL ser responsive, adaptándose correctamente a resoluciones de desktop (1024px+) y tablet (768px+).
4. WHEN el usuario redimensiona la ventana, THE Workspace SHALL ajustar las proporciones de los paneles manteniendo la usabilidad.

### Requirement 4: Carga de Datos del Ejercicio

**User Story:** As a estudiante, I want que al entrar a un nivel el editor tenga un código inicial (boilerplate), so that no empiece con el archivo en blanco y pueda enfocarme en la lógica.

#### Acceptance Criteria

1. WHEN el estudiante selecciona un Node, THE Frontend SHALL consumir el endpoint GET /nodes/{node_id}.
2. WHEN el Frontend recibe la respuesta del ejercicio, THE Frontend SHALL renderizar el enunciado (markdown_content) en el panel izquierdo con soporte Markdown básico.
3. WHEN el Frontend recibe la respuesta del ejercicio, THE Frontend SHALL inyectar el boilerplate_code automáticamente en el Monaco_Editor.
4. IF el endpoint retorna error, THEN THE Frontend SHALL mostrar un mensaje de error amigable al usuario.
5. WHEN el ejercicio se carga exitosamente, THE Frontend SHALL mostrar el título del Node en la cabecera del Workspace.
