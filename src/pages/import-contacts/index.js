// @flow

import React from "react";
import PageLayout from "../../style-guide/page-layout";
import { Display3 } from "baseui/typography";
import { Spacer20 } from "../../style-guide/spacer";
import UploadForm from "./upload-form";
import MatchColumns from "./match-columns";
import Success from "./success";
import type { CreateContactResult } from "../../lib/types/audience";

const Steps = {
  UploadFile: 1,
  MatchColumns: 2,
  Success: 3,
};

const ImportContactsPage = () => {
  const [step, setStep] = React.useState(Steps.UploadFile);
  const [csvFile, setCSVFile] = React.useState<?File>(null);
  const [uploadResults, setUploadResults] = React.useState<
    CreateContactResult[]
  >([]);

  return (
    <PageLayout>
      <Display3>Import Contacts</Display3>
      <Spacer20 />
      {step === Steps.UploadFile && (
        <UploadForm
          onUpload={(file) => {
            setCSVFile(file);
            setStep(Steps.MatchColumns);
          }}
        />
      )}

      {step === Steps.MatchColumns && csvFile ? (
        <MatchColumns
          csvFile={csvFile}
          onBack={() => setStep(Steps.UploadFile)}
          onUpload={(results) => {
            setUploadResults(results);
            setStep(Steps.Success);
          }}
        />
      ) : null}

      {step === Steps.Success && <Success results={uploadResults} />}
    </PageLayout>
  );
};

export default ImportContactsPage;
