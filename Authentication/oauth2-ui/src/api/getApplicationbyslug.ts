import { authenticationPrivateApi } from "../core/config";
export async function getApplicationbyslug(slug: string) {
    return await authenticationPrivateApi.get(`/api/authentication/v1/config/${slug}`).then(({ data }) => data);
}
