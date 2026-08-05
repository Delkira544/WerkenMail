import "@typespec/http";

using Http;

model FieldError {
  field: string;
  message: string;
}

/** Códigos de error del wraper en shared/errors (errors.go). */
enum ErrorCode {
  INTERNAL_ERROR,
  UNAUTHORIZED,
  NOT_FOUND,
  FORBIDDEN,
  BAD_REQUEST,
  CONFLICT,
  VALIDATION_ERROR,
}

/** Cuerpo del error (response.APIError). */
model ApiError {
  code: ErrorCode;
  message: string;
  details?: FieldError[];
}

/** Envelope de error: success=false + error. */
@error
model ErrorEnvelope {
  success: false;
  error: ApiError;
}

/** Metadatos de paginación (response.Meta). */
model Meta {
  request_id: string;
  timestamp: string;
  page?: int32;
  per_page?: int32;
  total?: int64;
  total_pages?: int32;
  has_next?: boolean;
  has_prev?: boolean;
}

/** Envelope de éxito: success=true + data (+ meta opcional). */
model SuccessEnvelope<T> {
  success: true;
  data: T;
}

model SuccessWithMetaEnvelope<T> {
  success: true;
  data: T;
  meta: Meta;
}

/** Envelope de éxito para HTTP 201 Created. */
model CreatedEnvelope<T> {
  @statusCode statusCode: 201;
  ...SuccessEnvelope<T>;
}
