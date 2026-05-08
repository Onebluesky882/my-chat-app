import React, { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Image,
} from "react-native";
import { ChevronLeft } from "lucide-react-native";
import { authClient } from "@/lib/auth-client";
import { SafeAreaView } from "react-native-safe-area-context";
export default function AuthScreen() {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [repassword, setRePassword] = useState("");

  const handleAuth = async () => {
    if (isLogin) {
      await authClient.signIn.email({ email, password });
    } else {
      if (password !== repassword) return alert("Passwords do not match");
      await authClient.signUp.email({ email, password, name });
    }
  };

  return (
    // bg-[#071219] คือสีพื้นหลังน้ำเงินเข้มตามรูป Messenger
    <SafeAreaView className="flex-1 bg-[#071219]">
      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        className="flex-1"
      >
        <ScrollView contentContainerStyle={{ flexGrow: 1 }} className="px-6">
          {/* Header: Back Button */}
          <TouchableOpacity className="mt-4 -ml-2">
            <ChevronLeft color="white" size={30} />
          </TouchableOpacity>

          {/* Logo Section */}
          <View className="items-center mt-10 mb-12">
            <Image
              source={{
                uri: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/be/Facebook_Messenger_logo_2020.svg/512px-Facebook_Messenger_logo_2020.svg.png",
              }}
              className="w-20 h-20"
            />
          </View>

          {/* Input Form */}
          <View className="w-full">
            {!isLogin && (
              <TextInput
                className="bg-[#1C2733] border border-[#2D3A4A] text-white rounded-2xl px-4 py-4 text-base mb-3"
                placeholder="Full Name"
                placeholderTextColor="#888"
                value={name}
                onChangeText={setName}
              />
            )}

            <TextInput
              className="bg-[#1C2733] border border-[#2D3A4A] text-white rounded-2xl px-4 py-4 text-base mb-3"
              placeholder="Mobile number or email"
              placeholderTextColor="#888"
              value={email}
              onChangeText={setEmail}
              autoCapitalize="none"
            />

            <TextInput
              className="bg-[#1C2733] border border-[#2D3A4A] text-white rounded-2xl px-4 py-4 text-base mb-3"
              placeholder="Password"
              placeholderTextColor="#888"
              value={password}
              onChangeText={setPassword}
              secureTextEntry
            />

            {!isLogin && (
              <TextInput
                className="bg-[#1C2733] border border-[#2D3A4A] text-white rounded-2xl px-4 py-4 text-base mb-3"
                placeholder="Confirm Password"
                placeholderTextColor="#888"
                value={repassword}
                onChangeText={setRePassword}
                secureTextEntry
              />
            )}

            {/* Main Action Button */}
            <TouchableOpacity
              className="bg-[#0064E0] rounded-full py-4 items-center mt-2 shadow-lg"
              onPress={handleAuth}
            >
              <Text className="text-white font-semibold text-lg">
                {isLogin ? "Log in" : "Create Account"}
              </Text>
            </TouchableOpacity>

            <TouchableOpacity className="mt-5">
              <Text className="text-white text-center font-medium">
                Forgot password?
              </Text>
            </TouchableOpacity>
          </View>

          {/* Footer Section */}
          <View className="flex-1 justify-end items-center pb-8 mt-10">
            <TouchableOpacity
              className="w-full border border-[#0064E0] rounded-full py-3 mb-6"
              onPress={() => setIsLogin(!isLogin)}
            >
              <Text className="text-[#0064E0] text-center font-bold">
                {isLogin
                  ? "Create new account"
                  : "Already have an account? Log in"}
              </Text>
            </TouchableOpacity>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}
