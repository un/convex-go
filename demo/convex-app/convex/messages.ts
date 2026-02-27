import {
  action,
  mutation,
  query
} from "./_generated/server";
import { internal } from "./_generated/api";
import { v } from "convex/values";

const messageValidator = v.object({
  _id: v.id("messages"),
  _creationTime: v.number(),
  text: v.string(),
  source: v.union(
    v.literal("seed"),
    v.literal("mutation"),
    v.literal("action"),
    v.literal("dashboard")
  ),
  updatedAt: v.number()
});

export const list = query({
  args: {},
  returns: v.array(messageValidator),
  handler: async (ctx) => {
    return await ctx.db.query("messages").withIndex("by_updatedAt").order("desc").take(50);
  }
});

export const listWithLimit = query({
  args: {
    limit: v.number()
  },
  returns: v.array(messageValidator),
  handler: async (ctx, args) => {
    const limit = Math.min(Math.max(args.limit, 1), 100);
    return await ctx.db.query("messages").withIndex("by_updatedAt").order("desc").take(limit);
  }
});

export const create = mutation({
  args: {
    text: v.optional(v.string()),
    source: v.optional(
      v.union(v.literal("seed"), v.literal("mutation"), v.literal("action"), v.literal("dashboard"))
    )
  },
  returns: v.object({
    id: v.id("messages"),
    text: v.string(),
    source: v.union(v.literal("seed"), v.literal("mutation"), v.literal("action"), v.literal("dashboard"))
  }),
  handler: async (ctx, args) => {
    const now = Date.now();
    const text = args.text ?? `Message created at ${new Date(now).toISOString()}`;
    const source = args.source ?? "mutation";
    const id = await ctx.db.insert("messages", {
      text,
      source,
      updatedAt: now
    });
    return { id, text, source };
  }
});

export const updateText = mutation({
  args: {
    id: v.id("messages"),
    text: v.string()
  },
  returns: v.null(),
  handler: async (ctx, args) => {
    await ctx.db.patch(args.id, {
      text: args.text,
      source: "dashboard",
      updatedAt: Date.now()
    });
    return null;
  }
});

export const seedIfEmpty = mutation({
  args: {},
  returns: v.object({
    inserted: v.number()
  }),
  handler: async (ctx) => {
    const existing = await ctx.db.query("messages").take(1);
    if (existing.length > 0) {
      return { inserted: 0 };
    }

    const now = Date.now();
    await ctx.db.insert("messages", {
      text: "Seed message: open /api/live/stream and edit this row in dashboard",
      source: "seed",
      updatedAt: now
    });
    await ctx.db.insert("messages", {
      text: "Seed message: realtime updates should appear in the stream",
      source: "seed",
      updatedAt: now + 1
    });
    return { inserted: 2 };
  }
});

export const ping = action({
  args: {
    note: v.optional(v.string())
  },
  returns: v.object({
    ok: v.boolean(),
    note: v.string(),
    messageCount: v.number(),
    latestText: v.union(v.string(), v.null())
  }),
  handler: async (ctx, args): Promise<{
    ok: boolean;
    note: string;
    messageCount: number;
    latestText: string | null;
  }> => {
    const note = args.note ?? "ping from action";
    await ctx.runMutation(internal.messages_internal.insertActionLog, { note });
    const summary = await ctx.runQuery(internal.messages_internal.summary, {});
    return {
      ok: true,
      note,
      messageCount: summary.count,
      latestText: summary.latestText
    };
  }
});
