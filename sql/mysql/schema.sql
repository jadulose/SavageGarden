create table class
(
    c_id   int auto_increment
        primary key,
    c_name varchar(20) not null,
    c_s_id int         null
);

create table student
(
    s_id   int         not null
        primary key,
    s_name varchar(20) not null,
    s_mail varchar(50) not null,
    s_pass binary(60)  null,
    s_c_id int         null,
    constraint student_class_cid_fk
        foreign key (s_c_id) references class (c_id)
            on update cascade on delete cascade
);

alter table class
    add constraint class_student_sid_fk
        foreign key (c_s_id) references student (s_id)
            on update cascade on delete set null;

create table subject
(
    sj_id    int auto_increment
        primary key,
    sj_name  varchar(20) not null,
    sj_tname varchar(20) not null,
    sj_tmail varchar(50) not null
);

create table homework
(
    h_id       int auto_increment
        primary key,
    h_title    varchar(30)  not null,
    h_detail   varchar(500) not null,
    h_deadline datetime     not null,
    h_sj_id    int          not null,
    constraint homework_subject_sj_id_fk
        foreign key (h_sj_id) references subject (sj_id)
            on update cascade on delete cascade
);

create table submit
(
    sm_id      int auto_increment
        primary key,
    sm_h_id    int                                    not null,
    sm_s_id    int                                    not null,
    sm_name    varchar(50)                            not null,
    sm_time    datetime   default current_timestamp() not null,
    sm_ip      char(15)                               not null,
    sm_os      varchar(10)                            not null,
    sm_browser varchar(20)                            not null,
    sm_login   tinyint(1) default 0                   not null,
    constraint submit_homework_h_id_fk
        foreign key (sm_h_id) references homework (h_id)
            on update cascade on delete cascade,
    constraint submit_student_s_id_fk
        foreign key (sm_s_id) references student (s_id)
            on update cascade on delete cascade
);

create table choose
(
    ch_h_id  int not null,
    ch_s_id  int not null,
    ch_sm_id int not null,
    primary key (ch_h_id, ch_s_id),
    constraint choose_homework_h_id_fk
        foreign key (ch_h_id) references homework (h_id),
    constraint choose_student_s_id_fk
        foreign key (ch_s_id) references student (s_id),
    constraint choose_submit_sm_id_fk
        foreign key (ch_sm_id) references submit (sm_id)
);

