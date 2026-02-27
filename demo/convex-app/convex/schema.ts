import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  messages: defineTable({
    text: v.string(),
    source: v.union(
      v.literal("seed"),
      v.literal("mutation"),
      v.literal("action"),
      v.literal("dashboard")
    ),
    updatedAt: v.number()
  })
    .index("by_updatedAt", ["updatedAt"])
    .index("by_source", ["source"])
});
