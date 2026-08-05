import "@typespec/http";
import "../common.tsp";

using Http;

namespace WerkenMailApi.Auth;

model LoginRequest {
  username: string;
  password: string;
}

model LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_at: int64;
  role: string;
}

@route("/auth")
@tag("Auth")
interface AuthService {
  /** Autentica un usuario y devuelve la sesión. */
  @useAuth(NoAuth)
  @route("/login")
  @summary("Autenticate")
  @doc("Authenticate a user and return the session.")
  @post
  login(@body body: LoginRequest):
    | SuccessEnvelope<LoginResponse>
    | ErrorEnvelope;
}
