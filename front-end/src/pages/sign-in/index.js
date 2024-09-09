import { useState, useEffect } from "react";
import Link from "next/link";
import Head from "next/head";
import { useRouter } from "next/router";

import styles from "./styles.module.css";
import handleSignInAPI from "../../apis/handlers/sign-in";

import { Layout, theme, Form, Input, Button, message } from "antd";
import { GetToken } from "../../components/jwt/index";

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
  const router = useRouter();

  useEffect(() => {
    const token = GetToken();

    if (token) {
      router.push("/chat");
      return;
    }
  }, []);

  const {
    token: { colorBgContainer },
  } = theme.useToken();

  const handleSignIn = (e) => {
    e.preventDefault();
    if (!credentialsValidate(username, password)) {
      handleSignInAPI(username, password)
        .then((data) => {
          localStorage.setItem("token", data.token);
          router.push("/chat");
        })
        .catch((error) => {
          console.log("Error getting delivery address: ", error);
        });
    }
  };

  return (
    <Layout>
      <Head>
        <title>Sign In</title>
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
