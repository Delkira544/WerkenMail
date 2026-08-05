# WerkenMail

**WerkenMail** es un cliente SMTP administrado que permite a equipos y proyectos de desarrollo enviar mensajes de forma centralizada a través de una API REST, sin necesidad de que cada proyecto configure su propio proveedor de correo.

## Descripción general

La plataforma permite crear **proyectos**, cada uno con sus propias **API Keys** con alcance (scope) limitado, que se utilizan para autenticar las solicitudes contra la API de mensajería. Los proyectos comparten un único cliente SMTP ya configurado a nivel de infraestructura, lo que simplifica la administración de dominios y remitentes verificados.

El sistema incorpora un **gestor de plantillas de correo** (similar al enfoque de Resend/Postmark/SendGrid), donde cada plantilla se define en HTML con variables dinámicas mediante la sintaxis `{{variable}}`. Al momento del envío, el cliente solo necesita indicar la plantilla a utilizar y los valores de las variables correspondientes; el backend se encarga de renderizar el mensaje final.

## Autenticación

WerkenMail utiliza un esquema de autenticación dual:

- **JWT (derivado de LDAP)**: utilizado para el acceso al dashboard administrativo. Los roles se resuelven una única vez durante el login a partir del `gidNumber` de LDAP (por ejemplo, `600` para roles administrativos/funcionales y `500` para roles de estudiante), y se incluyen como claims en el token.
- **API Keys**: utilizadas para autenticar las solicitudes programáticas contra la API de mensajería, con alcance limitado por proyecto.

## Características principales

- Gestión de **proyectos** con miembros y roles asociados.
- **Plantillas de correo** almacenadas en base de datos, con extracción automática de variables y versionado.
- Renderizado de mensajes al momento del envío, preservando un snapshot inmutable del contenido y las variables utilizadas en cada mensaje enviado.
- **API Keys** con prefijo visible y hash seguro, con revelación única del valor en texto plano al momento de su creación.
- Sanitización de contenido HTML dinámico antes de su interpolación en las plantillas, garantizando compatibilidad con los principales clientes de correo (Gmail, Outlook, etc.).
- Envío asíncrono de mensajes mediante un único cliente SMTP compartido a nivel de infraestructura.

## Objetivo

Ofrecer a equipos de desarrollo una solución simple y centralizada para el envío transaccional de correos desde sus propias aplicaciones o servicios, evitando la necesidad de administrar credenciales SMTP por proyecto y estandarizando el uso de plantillas dinámicas mediante una API REST.
