import os
import csv
import re
import xlrd
from flask import request, after_this_request, make_response, send_file
from collections import OrderedDict


# For handling csv errors
class InvalidUsage(Exception):
    status_code = 400

    def __init__(self, message, status_code=None, payload=None):
        Exception.__init__(self)
        self.message = message
        if status_code is not None:
            self.status_code = status_code
        self.payload = payload

    def to_dict(self):
        rv = dict(self.payload or ())
        rv['message'] = self.message
        return rv


# error handler for invalid csv
def handle_error(errMsg, status_code=400):
    errMsg += '\nPlease fix error and try uploading again.'
    resp = make_response(errMsg, status_code)
    return resp


# inverse SKU map for deSKU
inverseSKU = {
    "A1": "Search",
    "A2": "Display",
    "A3": "Social",
    "A4": "Email",
    "A5": "Ad_Video",
    "B1": "KWD",
    "B2": "MP",
    "B3": "MP-KWD",
    "B4": "Topic",
    "B5": "Topic-KWD",
    "B6": "Connection",
    "B7": "Behavior",
    "B8": "Demographic",
    "B9": "Interest",
    "B10": "Lookalike",
    "B11": "Remarketing",
    "B12": "Newsletter - About",
    "B13": "Newsletter - Programs",
    "B14": "Newsletter - Addiction",
    "B15": "Newsletter - PTSD",
    "B16": "Newsletter - MH",
    "B17": "Run of Site",
    "B18": "Geographic",
    "B19": "Conversion",
    "B20": "Newsletter - Mood",
    "B21": "No Response - SF",
    "B22": "Custom Segments",
    "B23": "Dynamic - Domain",
    "B24": "Dynamic - SSE",
    "B25": "Question",
    "B26": "Question Retargeting",
    "B27": "Sign-Up",
    "B28": "Eblast",
    "B29": "No Response",
    "B30": "Smart Display",
    "C1": "Canvas Ad",
    "C2": "2 Headline",
    "C3": "GIF",
    "C4": "Gmail Ad",
    "C5": "HTML",
    "C6": "Image Carousel",
    "C7": "Lightbox",
    "C8": "Link Post",
    "C9": "Photo Post",
    "C10": "Responsive",
    "C11": "Static Image",
    "C12": "Video",
    "C13": "Video Carousel",
    "C14": "Long Content - No Image",
    "C15": "Long Content - Image",
    "C16": "Short Content - No Image",
    "C17": "Short Content - Image",
    "C18": "Banner",
    "C19": "Profile",
    "C20": "Text Ad",
    "C21": "3 Headline",
    "C22": "Dynamic Headline",
    "C23": "Native",
    "D1": "About-General",
    "D2": "About-Other",
    "D3": "About-Self",
    "D4": "Benefits-General",
    "D5": "Benefits-Other",
    "D6": "Benefits-Self",
    "D7": "Emotion-General",
    "D8": "Emotion-Other",
    "D9": "Emotion-Self",
    "D10": "Leading-General",
    "D11": "Leading-Other",
    "D12": "Leading-Self",
    "D13": "Scare-General",
    "D14": "Scare-Other",
    "D15": "Scare-Self",
    "D16": "Stats-General",
    "D17": "Stats-Other",
    "D18": "Stats-Self",
    "D19": "Urgent-General",
    "D20": "Urgent-Other",
    "D21": "Urgent-Self",
    "D22": "Dynamic Message",
    "D23": "Info-General",
    "D24": "Info-Other",
    "D25": "Info-Self",
    "D26": "Question-General",
    "D27": "Question-Other",
    "D28": "Question-Self",
    "D29": "Hope-General",
    "D30": "Hope-Other",
    "D31": "Hope-Self",
    "E1": "18-24",
    "E2": "25-34",
    "E3": "35-44",
    "E4": "35-49",
    "E5": "45-54",
    "E6": "50-64",
    "E7": "55-64",
    "E8": "65+",
    "E9": "All",
    "E10": "Undetermined",
    "F1": "African American",
    "F2": "Asian",
    "F3": "Hispanic",
    "F4": "Native American",
    "F5": "Pacific Islander",
    "F6": "Two or More Races",
    "F7": "White",
    "F8": "Undetermined",
    "F9": "All",
    "G1": "Brother",
    "G2": "Daughter",
    "G3": "Husband",
    "G4": "Parent - Expecting",
    "G5": "Parent - 0-12 Month",
    "G6": "Parent - Pre-Teen",
    "G7": "Parent - Teen",
    "G8": "Parent - Adult Child",
    "G9": "Sister",
    "G10": "Son",
    "G11": "Wife",
    "G12": "Undetermined",
    "G13": "All",
    "G14": "Parent",
    "H1": "Female",
    "H2": "Male",
    "H3": "Undetermined",
    "H4": "All",
    "I1": "39k or Less",
    "I2": "40k-49k",
    "I3": "50k-74k",
    "I4": "75k-99k",
    "I5": "100k-124k",
    "I6": "125k-149k",
    "I7": "150k-249k",
    "I8": "250k-349k",
    "I9": "350k-499k",
    "I10": "500k or More",
    "I11": "Undetermined",
    "I12": "All",
    "J1": "Away from Family",
    "J2": "Away from Home",
    "J3": "Business and Industry",
    "J4": "Chronic Relapser",
    "J5": "Democrat",
    "J6": "Detox Seeker",
    "J7": "Entertainment",
    "J8": "Fitness and Wellness",
    "J9": "Food and Drink",
    "J10": "Friends of Alumni",
    "J11": "LGBT Population",
    "J12": "Outdoors",
    "J13": "Politics",
    "J14": "Previous Patient of Competitor",
    "J15": "Republican",
    "J16": "Shopping and Fashion",
    "J17": "Sports",
    "J18": "Technology",
    "J19": "Travel",
    "J20": "Undetermined",
    "J21": "All",
    "J22": "Employers",
    "K1": "English",
    "K2": "Spanish",
    "K3": "Bilingual",
    "K4": "Undetermined",
    "K5": "All",
    "L1": "Some High School",
    "L2": "High School Grad",
    "L3": "Associate Degree",
    "L4": "Some College",
    "L5": "College Grad",
    "L6": "Professional Degree",
    "L7": "Some Grad School",
    "L8": "Master's Degree",
    "L9": "Doctorate Degree",
    "L10": "Undetermined",
    "L11": "All",
    "M1": "Admin",
    "M2": "Arts",
    "M3": "Business and Finance",
    "M4": "Executive",
    "M5": "Government",
    "M6": "Healthcare",
    "M7": "IT",
    "M8": "Legal",
    "M9": "Manufacturing",
    "M10": "Sales",
    "M11": "Service",
    "M12": "Student-College",
    "M13": "Student-Grad School",
    "M14": "Student-High School",
    "M15": "Unemployed",
    "M16": "Undetermined",
    "M17": "All",
    "N1": "Divorced",
    "N2": "Married",
    "N3": "Separated",
    "N4": "Single",
    "N5": "Undetermined",
    "N6": "All",
    "O1": "Agnosticism",
    "O2": "Atheism",
    "O3": "Buddhism",
    "O4": "Christianity",
    "O5": "Hindu",
    "O6": "Islam",
    "O7": "Judaism",
    "O8": "Mormonism",
    "O9": "Sikhism",
    "O10": "Undetermined",
    "O11": "All"
}


# receive filepath/name and process csv; handing it back to route
def processDeSKU(filepath, reFilename, reFilepath, uploadFolder="/tmp/"):
    # check file extension
    _, file_extension = os.path.splitext(filepath)
    # column headers that we will insert into new csv
    headerKeys = ["Network", "Targeting_Method", "Format", "Message", "Age_Range", "Ethnicity", "Family_Role",
                  "Gender", "Income", "Interests/Behaviors", "Language", "Education", "Occupation",
                  "Relationship", "Religion"]

    # if it's an excel file
    if file_extension in ('.xlsm', '.xlsx', '.xls'):
        for ext in ('.xlsm', '.xlsx', '.xls'):
            reFilename = reFilename.replace(ext, '.csv')
        reFilepath = os.path.join(uploadFolder, reFilename)

        # create workbook
        wb = xlrd.open_workbook(filepath)
        worksheet = wb.sheet_by_index(0)
        hdrs = [str(cell.value) for cell in worksheet.row(0)]

        # start data list with file headers appended with new audit headers
        data_list = [hdrs + headerKeys]
        skuKey = hdrs[0]
        rowCount = 0
        for rowx in range(worksheet.nrows):
            if rowx < 1:  # (Optionally) skip headers
                continue

            rowCount += 1  # increment at beginning to account for header row
            # create dict row so we can append to data_list and write to new csv; same as we do in csv upload
            row = OrderedDict(zip(hdrs, worksheet.row_values(rowx)))

            # get sku from input file
            sku = row[skuKey]
            newRow = [sku]

            # explode the sku
            exploded = re.findall('[A-Z][^A-Z]*', sku)
            # exploded.pop(0)  # remove empty space at first index

            for i in range(0, len(headerKeys)):
                # iterate through headerkeys and matching exploded sku inverse values
                newRow.append(inverseSKU[exploded[i]])

            # add to temp data list after we explode and assign sku
            data_list.append(newRow)
        # END Loop

        # # write out to new csv file for response
        try:
            with open(reFilepath, "w") as fp:
                dict_writer = csv.writer(fp)
                dict_writer.writerows(data_list)
        except Exception as e:
            return handle_error("Write to CSV Err %s" % e)

    else:
        # parse it
        with open(filepath, "r") as f:
            reader = csv.DictReader(f)
            # skip header row before looping through to read
            hdrs = reader.fieldnames
            print(hdrs)
            # temp data list that we'll use to write to output csv file
            # start with headers from upload, then add ours
            data_list = [hdrs + headerKeys]
            skuKey = hdrs[0]
            for row in reader:
                # get sku from input file
                sku = row[skuKey]
                newRow = [sku]

                # explode the sku
                exploded = re.findall('[A-Z][^A-Z]*', sku)
                # exploded.pop(0)  # remove empty space at first index

                for i in range(0, len(headerKeys)):
                    # iterate through headerkeys and matching exploded sku inverse values
                    newRow.append(inverseSKU[exploded[i]])

                # add to temp data list after we explode and assign sku
                data_list.append(newRow)

        # write out to new file for response
        try:
            with open(reFilepath, "w") as fp:
                dict_writer = csv.writer(fp)
                dict_writer.writerows(data_list)
        except Exception as e:
            return handle_error("Write to CSV Err %s" % e)

    # delete the files once we're done; after file is returned
    @after_this_request
    def remove_file(response):
        try:
            os.remove(filepath)
            os.remove(reFilepath)
            return response
        except OSError:
            pass

    # return new file in response
    return send_file(reFilepath, attachment_filename=reFilename, as_attachment=True)
