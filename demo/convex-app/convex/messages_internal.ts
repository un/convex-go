import { internalAction, internalMutation, internalQuery } from "./_generated/server";
import { v } from "convex/values";

export const summary = internalQuery({
  args: {},
  returns: v.object({
    count: v.number(),
    latestText: v.union(v.string(), v.null())
  }),
  handler: async (ctx) => {
    const latest = await ctx.db.query("messages").withIndex("by_updatedAt").order("desc").take(1);
    const all = await ctx.db.query("messages").collect();
    return {
      count: all.length,
      latestText: latest[0]?.text ?? null
    };
  }
});

export const insertActionLog = internalMutation({
  args: {
    note: v.string()
  },
  returns: v.id("messages"),
  handler: async (ctx, args) => {
    return await ctx.db.insert("messages", {
      text: `[action] ${args.note}`,
      source: "action",
      updatedAt: Date.now()
    });
  }
});

export const echoAction = internalAction({
  args: {
    value: v.string()
  },
  returns: v.string(),
  handler: async (_ctx, args) => {
    return `echo:${args.value}`;
  }
});
