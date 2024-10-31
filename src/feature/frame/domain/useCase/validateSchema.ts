import Ajv from "ajv";

export async function validateJsonWithDynamicSchema(jsonData: any, schemaUrl: string) {
  const ajv = new Ajv({ allErrors: true });
  try {
    const schemaResponse = await fetch(schemaUrl);
    const schema = await schemaResponse.json();
    const validate = ajv.compile(schema);
    const isValid = validate(jsonData);
    if (!isValid) {
      return { valid: false, errors: validate.errors };
    }
    return { valid: true, errors: null };
  } catch (error: any) {
    return { valid: false, errors: [error.message] };
  }
}
