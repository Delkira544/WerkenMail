import "./lib/common.tsp";
import "./lib/features/auth.tsp";
import "./lib/features/projects.tsp";
import "./lib/features/templates.tsp";
import "./lib/features/health.tsp";

using Http;

/** API WerkenMail: envoltorio compartido en lib/common.tsp, features en lib/features/. */
@service(#{ title: "WerkenMail API" })
@server("https://kan.bn/api/v1", "Servidor de producción")
@server("http://localhost:8080/api/v1", "Servidor de desarrollo")
@useAuth(MyBearerAuth)
namespace WerkenMailApi;

@doc("Bearer authentication header of the form `Bearer <token>`, where `<token>` is your auth token.")
model MyBearerAuth {
  type: AuthType.http;
  scheme: "bearer";
}
