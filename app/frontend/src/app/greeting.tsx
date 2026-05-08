import { View, Text, Button } from "react-native";
import React from "react";
import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { useRouter } from "expo-router";

export default function Greething() {
  const router = useRouter();
  return (
    <ThemedView className="flex-1 justify-center items-center">
      <ThemedText>greeting</ThemedText>

      <Button title="signup" onPress={() => router.push("/")}></Button>
    </ThemedView>
  );
}
