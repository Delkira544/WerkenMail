import "@typespec/http";
import "../common.tsp";

using Http;

namespace WerkenMailApi.Projects;

model CreateProjectRequest {
  name: string;
  description?: string;
}

model UpdateProjectRequest {
  name?: string;
  description?: string;
}

model Project {
  id: string;
  user_id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

@route("/projects")
@tag("Projects")
interface ProjectsService {
  /** Crea un nuevo proyecto. */
  @post
  create(@body body: CreateProjectRequest):
    | CreatedEnvelope<Project>
    | ErrorEnvelope;

  /** Lista los proyectos del usuario. */
  @get
  list(@query search?: string, @query limit?: int32, @query offset?: int32):
    | SuccessEnvelope<Project[]>
    | ErrorEnvelope;

  /** Actualiza un proyecto existente. */
  @patch
  update(@path id: string, @body body: UpdateProjectRequest):
    | SuccessEnvelope<Project>
    | ErrorEnvelope;

  /** Elimina un proyecto. */
  @delete
  delete(@path id: string): void | ErrorEnvelope;
}
