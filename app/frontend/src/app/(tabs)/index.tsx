import React from "react";
import { SafeAreaView } from "react-native-safe-area-context";
import { ThemedView } from "@/components/themed-view";
import { ThemedText } from "@/components/themed-text";
import {
  background,
  border,
  cornerRadius,
  padding,
} from "@expo/ui/swift-ui/modifiers";
export default function HomeScreen() {
  return (
    <SafeAreaView style={{ backgroundColor: "transparent" }}>
      <ThemedView>
        <ThemedText> Message</ThemedText>
      </ThemedView>
    </SafeAreaView>
  );
}
