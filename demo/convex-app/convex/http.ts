import { httpRouter } from "convex/server";
import { httpAction } from "./_generated/server";
import { api } from "./_generated/api";

const http = httpRouter();

http.route({
  path: "/health",
  method: "GET",
  handler: httpAction(async () => {
    return new Response(JSON.stringify({ ok: true, service: "convex-go-demo-backend" }), {
      status: 200,
      headers: { "content-type": "application/json" }
    });
  })
});

http.route({
  path: "/seed",
  method: "POST",
  handler: httpAction(async (ctx) => {
    const result = await ctx.runMutation(api.messages.seedIfEmpty, {});
    return new Response(JSON.stringify(result), {
      status: 200,
      headers: { "content-type": "application/json" }
    });
  })
});

export default http;
