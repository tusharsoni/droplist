// @flow

import React from "react";
import PageLayout from "../../style-guide/page-layout";
import { Spacer40, Spacer8 } from "../../style-guide/spacer";
import AudienceTable from "./table";
import { Pagination } from "baseui/pagination";
import useFetch from "use-http";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import type { AudienceSummary, Contact } from "../../lib/types/audience";
import { Link, useHistory, useLocation } from "react-router-dom";
import { Display3, Label1 } from "baseui/typography";
import { useStyletron } from "baseui";
import { Button, KIND, SIZE } from "baseui/button";
import PeopleSearchSvg from "../../style-guide/illustrations/people-search";
import CreateContactButton from "./create-contact-button";

const CONTACTS_PER_PAGE = 20;

const AudiencePage = () => {
  const history = useHistory();
  const location = useLocation();
  const [css] = useStyletron();
  const queryParams = new URLSearchParams(location.search);
  const page = parseInt(queryParams.get("page"), 10) || 1;

  const {
    get: loadSummary,
    loading: summaryLoading,
    error: summaryError,
    data: summary,
  } = useFetch<AudienceSummary>("/audience/summary", {}, []);

  const {
    get: loadContacts,
    loading: contactsLoading,
    error: contactsError,
    data: contacts,
  } = useFetch<Contact[]>(
    `/audience/contacts?limit=${CONTACTS_PER_PAGE}&offset=${
      (page - 1) * CONTACTS_PER_PAGE
    }`,
    {},
    [page]
  );

  const loading = summaryLoading || contactsLoading;
  const error = summaryError || contactsError;

  if (loading) {
    return (
      <PageLayout>
        <Spinner />
      </PageLayout>
    );
  }

  if (error) {
    return <PageLayout>Failed to load this page. Please try again.</PageLayout>;
  }

  return (
    <PageLayout>
      {summary.TotalContacts === 0 && (
        <div>
          <Display3>Audience</Display3>
          <Spacer40 />
          <div
            className={css({
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
            })}
          >
            <PeopleSearchSvg height={250} />
            <Spacer40 />
            <Label1>You don't have any contacts yet</Label1>
            <Spacer8 />
            <div className={css({ display: "flex" })}>
              <CreateContactButton
                onCreate={() => {
                  loadContacts();
                  loadSummary();
                }}
              />
              <Spacer8 />
              <Link
                className={css({ textDecoration: "none" })}
                to={"/audience/contacts/import"}
              >
                <Button kind={KIND.secondary} size={SIZE.compact}>
                  Import Contacts
                </Button>
              </Link>
            </div>
          </div>
        </div>
      )}

      {summary.TotalContacts > 0 && (
        <AudienceTable
          summary={summary}
          contacts={contacts}
          onRefresh={() => {
            loadSummary();
            loadContacts();
          }}
        />
      )}
      {summary.TotalContacts > CONTACTS_PER_PAGE && (
        <>
          <Spacer40 />
          <Pagination
            numPages={Math.ceil(summary.TotalContacts / CONTACTS_PER_PAGE)}
            currentPage={page}
            onPageChange={({ nextPage }) =>
              history.push(`/audience?page=${nextPage}`)
            }
          />
        </>
      )}
    </PageLayout>
  );
};

export default AudiencePage;
