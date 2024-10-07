import { useState, useEffect } from "react";
import Link from "next/link";
import Head from "next/head";
import Script from "next/script";
import { useRouter } from "next/router";

import styles from "./styles.module.css";
import handleSignInAPI from "../../apis/handlers/sign-in";
import handleValidateGoogleToken from "../../apis/handlers/google-sign-in";

import { Layout, theme, Form, Input, Button, message } from "antd";
import { GetToken } from "../../components/jwt/index";

require("dotenv").config();

const GOOGLE_CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;

const { Content } = Layout;

const credentialsValidate = (username, password) => {
  const formValidate = () => {
    if (!username) {
      return "Username must not be empty! Please try again";
    } else if (!password) {
      return "Password must not be empty! Please try again";
    }
    return null;
  };

  const validationError = formValidate();
  if (validationError) {
    return message.error(validationError);
  }
};

export default function SignInPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isClient, setIsClient] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const token = GetToken();
    if (token) {
      router.push("/chat");
      return;
    }

    setIsClient(true);
  }, []);

  const {
    token: { colorBgContainer },
  } = theme.useToken();

  // const handleGoogleSignIn = (googleUser) => {
  //   const id_token = googleUser.getAuthResponse().id_token;
  //   handleValidateGoogleToken(id_token)
  //     .then((data) => {
  //       localStorage.setItem("token", data.token);
  //       router.push("/chat");
  //     })
  //     .catch((error) => {
  //       message.error("Google Sign-In failed. Please try again.");
  //       console.error(error);
  //     });
  // };

  const handleSignIn = (e) => {
    e.preventDefault();
    if (!credentialsValidate(username, password)) {
      handleSignInAPI(username, password)
        .then((data) => {
          localStorage.setItem("token", data.token);
          router.push("/chat");
        })
        .catch((error) => {
          console.log("Error signing in: ", error);
          message.error("Sign-In failed. Please try again.");
        });
    }
  };

  return (
    <Layout>
      <Head>
        <title>Sign In</title>
        <meta charSet="UTF-8"></meta>
      </Head>

      <Layout
        style={{
          background: colorBgContainer,
          padding: 24,
          margin: 0,
          minHeight: 280,
        }}
      >
        <Content className={styles.content}>
          <Script
            src="https://accounts.google.com/gsi/client"
            async
            defer
          ></Script>
          <div
            id="g_id_onload"
            data-client_id={GOOGLE_CLIENT_ID}
            data-callback="handleCredentialResponse"
          ></div>
          <Form
            labelCol={{ span: 6 }}
            className={styles.signInForm}
            autoComplete="off"
          >
            <p className={styles.signInTitle}>Sign in</p>
            <Form.Item label="Username">
              <Input
                type="email"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </Form.Item>

            <Form.Item label="Password">
              <Input.Password
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </Form.Item>

            <p className={styles.forgotPassword}>Forgot Password</p>

            <Form.Item className={styles.signInButton}>
              <Button type="primary" htmlType="submit" onClick={handleSignIn}>
                Sign in
              </Button>
            </Form.Item>

            {isClient && <div class="g_id_signin" data-type="standard"></div>}

            <p className={styles.createAccount}>
              Don&#39;t have an account?{" "}
              <Link className={styles.signUp} href="/sign-up">
                Sign Up
              </Link>
            </p>
          </Form>
        </Content>
      </Layout>
    </Layout>
  );
}
