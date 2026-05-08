import "../global.css";

import {
  DarkTheme,
  DefaultTheme,
  ThemeProvider,
} from "@react-navigation/native";
import React from "react";
import { useColorScheme } from "react-native";

import { Stack } from "expo-router";
import { SafeAreaProvider } from "react-native-safe-area-context";

// mock
const isLoggedIn = false;

export default function TabLayout() {
  const colorScheme = useColorScheme();
  return (
    <ThemeProvider value={colorScheme === "dark" ? DarkTheme : DefaultTheme}>
      <SafeAreaProvider>
        <Stack>
          <Stack.Protected guard={isLoggedIn}>
            <Stack.Screen name="(tabs)" />
          </Stack.Protected>
          <Stack.Protected guard={!isLoggedIn}>
            <Stack.Screen name="greeting" />
          </Stack.Protected>
        </Stack>
      </SafeAreaProvider>
    </ThemeProvider>
  );
}
