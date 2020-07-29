// @flow
import * as React from "react";
import { Card, StyledAction, StyledBody } from "baseui/card";
import { Button } from "baseui/button";
import { useStyletron } from "baseui";
import LogoSvg from "../../style-guide/illustrations/logo";
import { Spacer20, Spacer8 } from "../../style-guide/spacer";
import { FormControl } from "baseui/form-control";
import { Input } from "baseui/input";
import useFetch from "use-http";
import type { AuthResponse } from "../../lib/types/auth";
import { useHistory } from "react-router-dom";
import { storeSession } from "../../lib/auth";
import { KIND as NotificationKind, Notification } from "baseui/notification";

const LoginPage = () => {
  const [css] = useStyletron();
  const signupAPI = useFetch<AuthResponse>(`/auth/email-otp/signup`);
  const loginAPI = useFetch<AuthResponse>(`/auth/email-otp/login`);
  const [showVerificationInput, setShowVerificationInput] = React.useState(
    false
  );
  const [email, setEmail] = React.useState("");
  const [verificationCode, setVerificationCode] = React.useState("");
  const history = useHistory();

  const onNext = async () => {
    const resp = await signupAPI.post({
      email,
    });

    if (!signupAPI.response.ok) {
      return;
    }

    if (resp.session_token) {
      storeSession(resp.user_uuid, resp.session_token);
      history.push("/");
      return;
    }

    setShowVerificationInput(true);
  };

  const onLogin = async () => {
    const resp = await loginAPI.post({
      email,
      verification_code: parseInt(verificationCode, 10),
    });

    if (loginAPI.response.ok) {
      storeSession(resp.user_uuid, resp.session_token);
      history.push("/");
    }
  };

  return (
    <div
      className={css({
        height: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      })}
    >
      <Card overrides={{ Root: { style: { width: "500px" } } }}>
        <StyledBody>
          <LogoSvg height={36} />
          <Spacer20 />
          <FormControl label="Email">
            <Input
              autoFocus
              placeholder="jane@doe.com"
              value={email}
              onChange={(e) => {
                setEmail(e.target.value);
                setShowVerificationInput(false);
              }}
            />
          </FormControl>
          {showVerificationInput && (
            <FormControl label="Verification Code">
              <Input
                autoFocus
                placeholder="0000"
                value={verificationCode}
                onChange={(e) => setVerificationCode(e.target.value)}
              />
            </FormControl>
          )}
          <Spacer20 />
        </StyledBody>
        <StyledAction>
          <Button
            isLoading={signupAPI.loading}
            onClick={showVerificationInput ? onLogin : onNext}
            overrides={{
              BaseButton: { style: { width: "100%" } },
            }}
          >
            {showVerificationInput ? "Login" : "Next"}
          </Button>

          {signupAPI.error ? (
            <>
              <Spacer8 />
              <Notification
                kind={NotificationKind.negative}
                overrides={{
                  Body: { style: { width: "auto" } },
                }}
              >
                Failed to sign up. Please try again
              </Notification>
            </>
          ) : loginAPI.error ? (
            <>
              <Spacer8 />
              <Notification
                kind={NotificationKind.negative}
                overrides={{
                  Body: { style: { width: "auto" } },
                }}
              >
                Check your verification code and try again
              </Notification>
            </>
          ) : loginAPI.response.ok ? null : signupAPI.response.ok ? (
            <>
              <Spacer8 />
              <Notification
                kind={NotificationKind.positive}
                overrides={{
                  Body: { style: { width: "auto" } },
                }}
              >
                A verification code has been sent to your email
              </Notification>
            </>
          ) : null}
        </StyledAction>
      </Card>
    </div>
  );
};

export default LoginPage;
