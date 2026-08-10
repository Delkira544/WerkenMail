import "@typespec/http";
import "../common.tsp";

using Http;

namespace WerkenMailApi.Templates;

@example(#{ key: "user_name", type: "string", default_value: "Juan Pérez" })
model CreateTemplateVariableRequest {
  key: string;
  type: string;
  default_value?: string;
}

@example(#{
  name: "Bienvenida Usuario Nuevo",
  subject: "¡Bienvenido a {{app_name}}, {{user_name}}!",
  body_html: "<table role=\"presentation\" width=\"100%\"><tr><td style=\"padding: 20px; font-family: Arial;\"><h1>Hola {{user_name}}</h1><p>Gracias por registrarte en {{app_name}}.</p></td></tr></table>",
  body_text: "Hola {{user_name}}, gracias por registrarte en {{app_name}}.",
  variables: #[
    #{ key: "user_name", type: "string", default_value: "Juan Pérez" },
    #{ key: "app_name", type: "string", default_value: "WerkenMail" }
  ],
})
model CreateTemplateRequest {
  name: string;
  subject: string;
  body_html: string;
  body_text?: string;
  variables?: CreateTemplateVariableRequest[];
}

model TemplateVariable {
  id: string;
  key: string;
  type: string;
  default_value?: string;
}

model Template {
  id: string;
  name: string;
  subject: string;
  body_html?: string;
  body_text?: string;
  version: uint16;
  variables: TemplateVariable[];
  created_at: string;
  updated_at: string;
}

model TemplateListItemResponse {
  id: string;
  name: string;
  subject: string;
  variable_count: int32;
  version: uint16;
  created_at: string;
  updated_at: string;
}

@route("/projects/{projectId}/templates")
@tag("Templates")
interface ProjectTemplatesService {
  /** Create a template inside a project */
  @post
  @summary("Create Template")
  create(
    @path projectId: string,
    @body body: CreateTemplateRequest,
  ): CreatedStandardResponse<Template>;

  /** Get all templates for a project */
  @get
  @summary("Get all Templates for a Project")
  list(
    @path projectId: string,
    @query search?: string,
    @query limit?: int32,
    @query offset?: int32,
  ): SuccessEnvelope<TemplateListItemResponse[]> | ErrorEnvelope;
}

@route("/templates")
@tag("Templates")
interface TemplatesService {
  /** Get a template by ID */
  @get
  @route("/{templateId}")
  @summary("Get Template by ID")
  get(@path templateId: string): SuccessEnvelope<Template> | ErrorEnvelope;

  /** Update a template by ID */
  @put
  @route("/{templateId}")
  @summary("Update Template by ID")
  update(@path templateId: string, @body body: CreateTemplateRequest):
    | SuccessEnvelope<Template>
    | ErrorEnvelope;

  /** Delete a template by ID */
  @delete
  @route("/{templateId}")
  @summary("Delete Template by ID")
  delete(@path templateId: string): void | ErrorEnvelope;
}
