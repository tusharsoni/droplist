// @flow

import React from "react";
import PageLayout from "../../style-guide/page-layout";
import AudienceTableHeader from "./table-header";
import { Spacer40 } from "../../style-guide/spacer";
import AudienceTable from "./table";
import { Pagination } from "baseui/pagination";
import useFetch from "use-http";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import type { AudienceSummary, Contact } from "../../lib/types/audience";
import { useHistory, useLocation } from "react-router-dom";

const CONTACTS_PER_PAGE = 20;

const AudiencePage = () => {
  const history = useHistory();
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const page = parseInt(queryParams.get("page"), 10) || 1;

  const {
    loading: summaryLoading,
    error: summaryError,
    data: summary,
  } = useFetch<AudienceSummary>("/audience/summary", {}, []);

  const {
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
      <AudienceTableHeader
        totalContacts={summary.TotalContacts}
        subscribedContacts={summary.SubscribedContacts}
      />
      <Spacer40 />
      <AudienceTable contacts={contacts} />
      <Spacer40 />
      <Pagination
        numPages={Math.ceil(summary.TotalContacts / CONTACTS_PER_PAGE)}
        currentPage={page}
        onPageChange={({ nextPage }) =>
          history.push(`/audience?page=${nextPage}`)
        }
      />
    </PageLayout>
  );
};

export default AudiencePage;
