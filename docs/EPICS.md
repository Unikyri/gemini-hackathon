Product Backlog - Gemini Coding Path (MVP)
==========================================

Este documento define las Épicas (Grandes funcionalidades) del proyecto basándose en la visión del Product Owner.

E01: Gestión de Identidad y Persistencia (User Core)
----------------------------------------------------

Descripción: Sistema que permite el acceso dual (Invitado vs Registrado) y la gestión del estado del usuario.

Alcance:

*   Permitir acceso como "Invitado" (Guest) para probar sin guardar historial a largo plazo.
    
*   Registro e Inicio de Sesión (Auth) para usuarios recurrentes.
    
*   Persistencia de datos: Guardar el progreso del mapa y el código "a medias" (drafts) si el usuario está logueado.
    
*   Si un invitado decide registrarse, migrar su progreso temporal a la cuenta nueva.
    

E02: El Arquitecto (AI Path Generator)
--------------------------------------

Descripción: El cerebro del sistema. Integración con Gemini 3 (Vertex AI) para interpretar la intención del usuario y generar la estructura del curso.

Alcance:

*   Análisis de Prompt: Detectar si el usuario quiere un lenguaje (Go), un rol (Backend) o una tecnología.
    
*   Generación Incremental: Crear lotes de **5 nodos** (ejercicios) a la vez. Al completar el 5º, se generan los siguientes 5 adaptados al rendimiento.
    
*   Generación de Contenido: Crear para cada nodo: Enunciado, Código Boilerplate, Documentación contextual y Unit Tests (Principales y Secundarios).
    

E03: El Mapa de Aprendizaje (Learning Path UI)
----------------------------------------------

Descripción: La interfaz visual que representa el progreso del estudiante (Estilo "Candy Crush").

Alcance:

*   Visualización de grafos/nodos conectados.
    
*   Estados de los nodos: Bloqueado (Gris), Disponible (Color/Pulsando), Completado (Estrellas/Check).
    
*   Navegación entre lotes de nodos.
    

E04: El Workspace de Código (The Dojo)
--------------------------------------

Descripción: El entorno de desarrollo integrado en el navegador donde ocurre la acción.

Alcance:

*   **Panel Izquierdo:** Enunciado del problema.
    
*   **Panel Central:** Editor de código (con Syntax Highlighting para Go).
    
*   **Panel Derecho (Pestañas Dinámicas):**
    
    1.  **Documentación:** Info oficial filtrada y adaptada por IA para el reto actual.
        
    2.  **Misiones Secundarias:** Lista de casos de uso extra (Edge cases) que otorgan XP adicional.
        
*   Botones de acción: "Ejecutar/Test" y "Guardar borrador".
    

E05: El Juez y Feedback (Code Runner & Feedback System)
-------------------------------------------------------

Descripción: Sistema de ejecución segura y devolución de respuestas al usuario.

Alcance:

*   Sandbox: Ejecución aislada de código Go (Docker/Judge0).
    
*   Validación: Ejecución de go test. Separar tests obligatorios (Pase de nivel) de los tests secundarios (Bonus).
    
*   **Feedback Adaptativo:**
    
    *   Niveles bajos: Traducir errores de compilador a lenguaje natural/amigable.
        
    *   Niveles altos: Mostrar stderr crudo para acostumbrar al usuario.
        
*   Prevención: Validaciones básicas de seguridad (ASVS) para evitar inyecciones malignas.
    

E06: Gamificación y Recompensas
-------------------------------

Descripción: Sistema de incentivos para mantener al usuario enganchado.

Alcance:

*   Cálculo de XP: Base por completar + Bonus por misiones secundarias.
    
*   Pantalla de "Victoria": Animación al pasar los tests con resumen de logros.
    
*   Visualización de progreso en el perfil del usuario.
    

**Notas Técnicas para el MVP:**

*   Lenguaje soportado: **Go**.
    
*   Base de datos: PostgreSQL.
    
*   IA: Google Vertex AI (Gemini 3).