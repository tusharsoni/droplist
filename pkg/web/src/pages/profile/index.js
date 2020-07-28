// @flow
import * as React from "react";
import PageLayout from "../../style-guide/page-layout";
import { Card, StyledAction, StyledBody } from "baseui/card";
import { Button } from "baseui/button";
import { FormControl } from "baseui/form-control";
import { Input } from "baseui/input";
import { Spacer20, Spacer8 } from "../../style-guide/spacer";
import { Select } from "baseui/select";
import { AWS_REGIONS } from "./aws-regions";
import useFetch from "use-http";
import type { Profile } from "../../lib/types/profile";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";

const ProfilePage = () => {
  const fetchProfileAPI = useFetch<Profile>(`/profile`);
  const updateProfileAPI = useFetch<Profile>(`/profile`);
  const [accessKeyID, setAccessKeyID] = React.useState("");
  const [secretAccessKey, setSecretAccessKey] = React.useState("");
  const [region, setRegion] = React.useState([]);

  React.useEffect(() => {
    async function loadData() {
      const profile: Profile = await fetchProfileAPI.get();

      if (fetchProfileAPI.response.ok) {
        setRegion(
          [AWS_REGIONS.find((r) => r.id === profile.AWSRegion)].filter(Boolean)
        );
        setAccessKeyID(profile.AWSAccessKeyID);
        setSecretAccessKey(profile.AWSSecretAccessKey);
      }
    }

    loadData();
  }, []); /* eslint-disable-line react-hooks/exhaustive-deps */

  if (fetchProfileAPI.loading) {
    return (
      <PageLayout>
        <Spinner />
      </PageLayout>
    );
  }

  if (fetchProfileAPI.error && fetchProfileAPI.response.status !== 404) {
    return (
      <PageLayout>
        <Notification kind={NotificationKind.negative}>
          Failed to this page
        </Notification>
      </PageLayout>
    );
  }

  return (
    <PageLayout>
      <Card
        title="AWS Settings"
        overrides={{ Root: { style: { maxWidth: "500px", margin: "0 auto" } } }}
      >
        <StyledBody>
          <Spacer20 />
          <FormControl label={() => "Access Key ID"}>
            <Input
              value={accessKeyID}
              onChange={(e) => setAccessKeyID(e.target.value)}
            />
          </FormControl>
          <FormControl label={() => "Secret Access Key"}>
            <Input
              type="password"
              value={secretAccessKey}
              onChange={(e) => setSecretAccessKey(e.target.value)}
              overrides={{
                MaskToggleButton: () => null,
              }}
            />
          </FormControl>
          <FormControl label={() => "Region"}>
            <Select
              options={AWS_REGIONS}
              clearable={false}
              value={region}
              onChange={({ value }) => setRegion(value)}
            />
          </FormControl>
          <Spacer20 />
        </StyledBody>
        <StyledAction>
          <Button
            isLoading={updateProfileAPI.loading}
            onClick={() => {
              updateProfileAPI.post({
                aws_region: region.length ? region[0].id : null,
                aws_access_key_id: accessKeyID || null,
                aws_secret_access_key:
                  fetchProfileAPI.data.AWSSecretAccessKey !== secretAccessKey &&
                  secretAccessKey
                    ? secretAccessKey
                    : null,
              });
            }}
            overrides={{ BaseButton: { style: { width: "100%" } } }}
          >
            Save Changes
          </Button>
          {updateProfileAPI.error && (
            <>
              <Spacer8 />
              <Notification
                kind={NotificationKind.negative}
                overrides={{
                  Body: { style: { width: "auto" } },
                }}
              >
                Failed to save changes
              </Notification>
            </>
          )}

          {updateProfileAPI.response.ok && (
            <>
              <Spacer8 />
              <Notification
                kind={NotificationKind.positive}
                overrides={{
                  Body: { style: { width: "auto" } },
                }}
              >
                Saved changes successfully
              </Notification>
            </>
          )}
        </StyledAction>
      </Card>
    </PageLayout>
  );
};

export default ProfilePage;
